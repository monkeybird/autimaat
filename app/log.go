// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	// LogFormat defines the date layout for log file names.
	LogFormat = "20060102"

	// LogPurgeTimeout defines the timeout after which the bot should
	// check for stale log files.
	LogPurgeTimeout = time.Hour * 24

	// LogRefreshTimeout determines how often we should check if a new
	// log file should be opened.
	LogRefreshTimeout = time.Minute

	// LogExpiration defines how old a log file should be, before it
	// is considered stale.
	LogExpiration = time.Hour * 24 * 7 * 2
)

// These defines some internal state.
var (
	logFile     *os.File
	startOnce   sync.Once
	stopOnce    sync.Once
	logPollQuit = make(chan struct{})
)

// InitLog initializes a new log file, if necessary. It then launches a
// background service which periodically checks if a new log file should
// be created. This happens according to a predefined timeout. Additionally,
// it will periodically purge stale log files from disk.
func InitLog(dir string) {
	startOnce.Do(func() {
		err := openLog(dir)
		if err != nil {
			log.Println("[app] Init log:", err)
			return
		}

		go logPoll(dir)
	})
}

// ShutdownLog shuts down the background log operations.
func ShutdownLog() {
	stopOnce.Do(func() {
		close(logPollQuit)
	})
}

// logPoll periodically purges stale log files and ensures logs are cycled
// after the appropriate timeout.
func logPoll(dir string) {
	// Do an initial purge of stale logs. This ensures that we
	// do not accumulate stale files if the PurgeTimeout below
	// is never triggered. Which might happen if the program is
	// shut down before the timeout occurs.
	err := purgeLogs(dir)

loopy:
	for err == nil {
		select {
		case <-logPollQuit:
			break loopy

		case <-time.After(LogRefreshTimeout):
			err = openLog(dir)

		case <-time.After(LogPurgeTimeout):
			err = purgeLogs(dir)
		}
	}

	if err != nil {
		log.Println("[app]", err)
	}

	// Clean up the existing log file.
	if logFile != nil {
		log.SetOutput(os.Stderr)
		logFile.Close()
		logFile = nil
	}
}

// openLog opens a new, or existing log file.
func openLog(dir string) error {
	// Ensure the log file directory exists.
	err := os.Mkdir(dir, 0700)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Determine the name of the new log file.
	timeStamp := time.Now().Format(LogFormat)
	file := fmt.Sprintf("%s.txt", timeStamp)
	file = filepath.Join(dir, file)

	// Exit if we're already using this file.
	if logFile != nil && logFile.Name() == file {
		if logFile.Name() == file {
			return nil
		}

		log.Println("[log] Opening new log file:", file)
	}

	// Create/open the new logfile.
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	// Set the new log output.
	log.SetOutput(fd)

	// Close the old log file and assign the new one.
	if logFile != nil {
		logFile.Close()
	}

	logFile = fd

	// Set the log prefix to include our process id.
	// This makes analyzing log data a little easier.
	log.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))
	return nil
}

// purgeLogs checks the given directory for files which are older than a
// predefined number of days. If found, the log file in question is deleted.
// This ensures we do not keep stale logs around unnecessarily.
func purgeLogs(dir string) error {
	log.Println("[log] Purging stale log files...")

	fd, err := os.Open(dir)
	if err != nil {
		return err
	}

	files, err := fd.Readdir(-1)
	fd.Close()

	if err != nil {
		return err
	}

	for _, file := range files {
		if time.Since(file.ModTime()) < LogExpiration {
			continue
		}

		path := filepath.Join(dir, file.Name())
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
