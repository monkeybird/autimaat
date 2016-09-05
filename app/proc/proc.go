// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
Package proc provides process forking and initialization functionality for
one or more IRC clients. This is intended to facilitate zero-downtime binary
upgrades by allowing the bot to fork itself and passing existing network
connections to the new child process.

A parent forks itself through the Fork() call. The child process then gets
access to the client connections through a call to InheritedFiles().
Before calling this, ensure that `flag.Parse()` has been called at least once.
*/
package proc

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

// connectionCount defines the number of connections passed into a forked
// process.
var connectionCount uint

func init() {
	flag.UintVar(&connectionCount, "fork", 0, "Number of inherited file descriptors")
}

// Wait polls for OS signals to either kill or fork this process.
//
// The signals it waits for are: SIGKILL, SIGINT, SIGTERM and SIGUSR1.
// The latter one being responsible for forking this process. The others
// are there so we may cleanly exit this process.
func Wait(argv []string, clients ...*os.File) {
	signals := make(chan os.Signal, 1)
	signal.Notify(
		signals,
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGUSR1,
	)

	log.Println("[proc] Waiting for signals...")
	for sig := range signals {
		log.Println("[proc] received signal:", sig)
		if sig != syscall.SIGUSR1 {
			return
		}

		log.Println("[proc] forking process...")
		err := doFork(argv, clients...)
		if err != nil {
			log.Println(err)
		}
	}
}

// Kill sends SIGINT to the current process. This can be used to cleanly
// break out of a signal polling loop from anywhere in the program.
func Kill() { syscall.Kill(os.Getpid(), syscall.SIGINT) }

// KillParent sends SIGINT to the parent process. This is intended to be
// called by a child after it has been forked and has re-initialized the
// inherited connections. The parent may now shut down.
func KillParent() { syscall.Kill(os.Getppid(), syscall.SIGINT) }

// Fork sends SIGUSR1 to the current process. This kickstarts the
// forking process.
func Fork() { syscall.Kill(os.Getpid(), syscall.SIGUSR1) }

// doFork forks the current process into a child process and passes the
// given client connections along to be inherited.
//
// The specified argv list define custom command line parameters which should
// be used in the invocation.
//
// The client list contains any file descriptors which should be inherited
// by the client.
//
// The forked process is called with the `-fork N` command line parameter.
// Where N is the number of file descriptors being passed along. This is
// used by the InheritedFiles() call below to rebuild the files.
func doFork(argv []string, clients ...*os.File) error {
	// Build the command line arguments for our child process.
	// This includes any custom arguments defined in the profile.
	args := append([]string{"-fork", strconv.Itoa(len(clients))}, argv...)

	// Initialize the command runner.
	cmd := exec.Command(os.Args[0], args...)
	cmd.ExtraFiles = make([]*os.File, len(clients))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = clients

	// Fork the process.
	return cmd.Start()
}

// InheritedFiles returns a list of N file descriptors inherited from a
// previous session through the Fork call.
//
// This function assumes that flag.Parse() has been called at least once
// already. The `-fork` flag has been registered during initialization of
// this package.
func InheritedFiles() []*os.File {
	if connectionCount == 0 {
		return nil
	}

	out := make([]*os.File, connectionCount)

	for i := range out {
		out[i] = os.NewFile(3+uintptr(i), "conn"+strconv.Itoa(i))
	}

	return out
}
