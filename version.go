// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"runtime"
)

// Application name and version constants.
const (
	AppName         = "autimaat"
	AppVersionMajor = 0
	AppVersionMinor = 22
)

// AppVersionRevision defines the application build number.
// This needs to be a string, as it is defined externally
// through a command line option:
//
//     go install -ldflags "-X main.AppVersionRevision=`date -u +%s`"
//
var AppVersionRevision string

func init() {
	// Make sure the revision has a sane value.
	if len(AppVersionRevision) == 0 {
		AppVersionRevision = "0"
	}
}

// Version returns the application version as a string.
func Version() string {
	return fmt.Sprintf("%s %d.%d.%s (Go runtime %s).\nCopyright (c) 2016, Jim Teeuwen.",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionRevision, runtime.Version())
}
