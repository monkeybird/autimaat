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
	"os"
	"strconv"
	"syscall"
)

// connectionCount defines the number of connections passed into a forked
// process.
var connectionCount uint

func init() {
	flag.UintVar(&connectionCount, "fork", 0, "Number of inherited file descriptors")
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
