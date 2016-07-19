// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"monkeybird/irc"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Parse command line arguments and load the bot profile.
	profile := parseArgs()

	// Initialize the log.
	initLog(profile.Root())

	// Write PID file. It may be needed by a process supervisor.
	writePid()

	// Create the bot and open the connection.
	bot := New(profile)
	err := bot.Run()

	if err != nil {
		log.Println(err)
	}
}

// initLog initializes the log file and any other propertiues
// it might need. Log files are created anew, once a day.
// They are named after the current (local) date and stored in
// the $PROFILE_ROOT/logs/ directory.
func initLog(root string) {
	// Ensure the log file directory exists.
	logDir := filepath.Join(root, "logs")
	err := os.Mkdir(logDir, 0700)

	if err != nil && !os.IsExist(err) {
		fmt.Fprintln(os.Stderr, "failed to create log file directory:", err)
		os.Exit(1)
	}

	// Set the log target to a file.
	timeStamp := time.Now().Format("20060102")
	logFile := fmt.Sprintf("%s.txt", timeStamp)
	logFile = filepath.Join(logDir, logFile)

	fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open log file:", err)
		os.Exit(1)
	}

	log.SetOutput(fd)

	// Set the log prefix to include our process id.
	// This makes analyzing log data a little easier.
	log.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))
}

// writePid writes a file with process' pid. This is used by supervisors
// like systemd to track the process state.
func writePid() {
	fd, err := os.Create("app.pid")
	if err != nil {
		log.Println("[bot] Create PID file:", err)
		return
	}

	fmt.Fprintf(fd, "%d", os.Getpid())
	fd.Close()
}

// parseArgs parses and validates command line arguments.
func parseArgs() irc.Profile {
	flag.Usage = func() {
		fmt.Println("usage:", os.Args[0], "[options] <profile directory>")
		flag.PrintDefaults()
	}

	newconf := flag.Bool("new", false, "Create a new, default configuration file and exit.")
	version := flag.Bool("version", false, "Display version information.")
	flag.Parse()

	if *version {
		fmt.Println(Version())
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Read and validate the profile root directory.
	root, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Set root as current working directory.
	err = os.Chdir(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create a new bot profile instance.
	profile := irc.NewProfile(root)

	// If applicable, save a new, default profile and exit.
	if *newconf {
		err := profile.Save()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("New configuration saved.")
		fmt.Println("Please edit it and relaunch the program.")
		os.Exit(0)
	}

	// Load an existing profile.
	err = profile.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return profile
}
