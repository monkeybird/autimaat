# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.


REVISION=github.com/monkeybird/autimaat/app.VersionRevision
LDFLAGS_DEBUG="-X $(REVISION)=`date -u +%s`"
LDFLAGS_RELEASE="-w -s -X $(REVISION)=`date -u +%s`"


# Rebuild the binary and install it locally.
# Fork the currently running process, if applicable.
all:
	go install -ldflags $(LDFLAGS_DEBUG)
	kill -USR1 `pidof autimaat`

