# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

REMOTE_OS=freebsd
REMOTE_ARCH=amd64
LDFLAGS_DEBUG="-X main.AppVersionRevision=`date -u +%s`"
LDFLAGS_RELEASE="-w -s -X main.AppVersionRevision=`date -u +%s`"

all: debug


# Rebuild the binary for the current platform and make sure its revision
# number is incremented. Tell the running bot process it is time to fork
# itself.
debug:
	go install -ldflags $(LDFLAGS_DEBUG)
	kill -USR1 `pidof autimaat`


# Rebuild the binary for the target platform and make sure its revision
# number is incremented. Additionally, strip it of debug symbols.
#
# Then upload the binary to a remote machine, through ssh.
# Instruct the remote bot process to fork itself.
# Finally, clean up the build files.
release:
	GOOS=$(REMOTE_OS) GOARCH=$(REMOTE_ARCH) go build -ldflags $(LDFLAGS_RELEASE)
	# Uploading yet to be implemented.
	go clean