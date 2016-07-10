# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

all: debug

debug:
	# Rebuild the binary and make sure its revision number is incremented.
	go install -ldflags "-X main.AppVersionRevision=`date -u +%s`"
	# Tell the running bot process it is time to fork itself.
	kill -USR1 `pidof autimaat`

release:
	