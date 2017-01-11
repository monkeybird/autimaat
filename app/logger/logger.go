// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package logger defines facilities to write bot data to log files,
// along with code which cycles log cycles and purges log files
// when needed.
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	// Format defines the date layout for log file names.
	Format = "20060102"

	// PurgeTimeout defines the timeout after which the bot should
	// check for stale log files.
	PurgeTimeout = time.Hour * 24

	// RefreshTimeout determines how often we should check if a new
	// log file should be opened.
	RefreshTimeout = time.Minute

	// Expiration defines how old a log file should be, before it
	// is considered stale.
	Expiration = time.Hour * 24 * 7 * 2
)

// These defines some internal state.
var (
	logFile     *os.File
	startOnce   sync.Once
	stopOnce    sync.Once
	logPollQuit = make(chan struct{})
)

// Init initializes a new log file, if necessary. It then launches a
// background service which periodically checks if a new log file should
// be created. This happens according to a predefined timeout. Additionally,
// it will periodically purge stale log files from disk.
func Init(dir string) {
	startOnce.Do(func() {
		err := openLog(dir)
		if err != nil {
			log.Println("[app] Init log:", err)
			return
		}

		go poll(dir)
	})
}

// Shutdown shuts down the background log operations.
func Shutdown() {
	stopOnce.Do(func() {
		close(logPollQuit)
	})
}

// poll periodically purges stale log files and ensures logs are cycled
// after the appropriate timeout.
func poll(dir string) {

	refresh := time.Tick(RefreshTimeout)
	purgeCheck := time.Tick(PurgeTimeout)
	var err error

loopy:
	for err == nil {
		select {
		case <-logPollQuit:
			break loopy
		case <-refresh:
			err = openLog(dir)
		case <-purgeCheck:
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
	timeStamp := time.Now().Format(Format)
	file := fmt.Sprintf("%s.txt", timeStamp)
	file = filepath.Join(dir, file)

	// Exit if we're already using this file.
	if logFile != nil && logFile.Name() == file {
		return nil
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
		if time.Since(file.ModTime()) < Expiration {
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
