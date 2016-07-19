# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

# It rebuild the bot binary for the current platform and make sure
# its revision number is incremented. Tell the running bot process
# it is time to fork itself.

LDFLAGS_DEBUG="-X main.AppVersionRevision=`date -u +%s`"
LDFLAGS_RELEASE="-w -s -X main.AppVersionRevision=`date -u +%s`"

all: debug

debug:
	go install -ldflags $(LDFLAGS_DEBUG)
	kill -USR1 `pidof autimaat`


release:
	go install -ldflags $(LDFLAGS_RELEASE)
	kill -USR1 `pidof autimaat`