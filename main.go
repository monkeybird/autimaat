// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/monkeybird/autimaat/app"
	"github.com/monkeybird/autimaat/irc"
)

func main() {
	// Parse command line arguments and load the bot profile.
	profile := parseArgs()

	// Write PID file. It may be needed by a process supervisor.
	writePid()

	// Create and run the bot.
	err := Run(profile)
	if err != nil {
		log.Fatal("[bot]", err)
	}
}

// writePid writes a file with process' pid. This is used by supervisors.
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
		fmt.Println(app.Version())
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
