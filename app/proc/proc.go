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
	"os"
	"syscall"
)

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
