# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

TARG_OS="freebsd"
TARG_ARCH="amd64"
PATH_FINAL=""

all: debug


# Rebuild the binary for the current platform and make sure its revision
# number is incremented. Tell the running bot process it is time to fork
# itself.
debug:
	go install -ldflags "-X main.AppVersionRevision=`date -u +%s`"
	kill -USR1 `pidof autimaat`


# Rebuild the binary for the target platform and make sure its revision
# number is incremented. Additionally, strip it of debug symbols.
#
# Then upload the binary to remote machine, through ssh.
# Finally, clean up the build files.
release:
	GOROOT_FINAL=$(PATH_FINAL) GOOS=$(TARG_OS) GOARCH=$(TARG_ARCH) \
		go build -ldflags "-w -s -X main.AppVersionRevision=`date -u +%s`"
	go clean