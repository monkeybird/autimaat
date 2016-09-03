// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package app defines some utility code for the compiled binary.
// This information sits in its own package, because it is used
// in multiple sub-packages.
package app

import (
	"fmt"
	"runtime"
)

// Application name and version constants.
const (
	Name         = "autimaat"
	VersionMajor = 1
	VersionMinor = 5
)

// VersionRevision defines the application build number.
// This needs to be a string, as it is defined externally
// through a command line option:
//
//     go install -ldflags "-X github.com/monkeybird/autimaat/app.VersionRevision=`date -u +%s`"
//
// Note that the entire import path must be specified for this to work.
//
var VersionRevision string

func init() {
	// Make sure the revision has a sane value.
	if len(VersionRevision) == 0 {
		VersionRevision = "0"
	}
}

// Version returns the application version as a string.
func Version() string {
	return fmt.Sprintf("%s %d.%d.%s (Go runtime %s)",
		Name, VersionMajor, VersionMinor, VersionRevision, runtime.Version())
}
