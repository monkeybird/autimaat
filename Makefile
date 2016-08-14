# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

## Before using the release build mode, ensure the following
## environment variables are defined:
##
##   AUTIMAAT_DIR     : Remote directory to install binary in.
##   AUTIMAAT_HOST    : Remote SSH host name.
##   AUTIMAAT_PORT    : Remote SSH port number.
##
##
## The remote host should have a shell script present at: $AUTIMAAT_DIR/restart.sh
## It should contain the following code:
##
##     #!/usr/bin/env sh
##     kill -USR1 `pidof autimaat`
##

LDFLAGS_DEBUG="-X main.AppVersionRevision=`date -u +%s`"
LDFLAGS_RELEASE="-w -s -X main.AppVersionRevision=`date -u +%s`"

all: debug

## Rebuild the binary and install it locally.
## Fork the currently running process, if applicable.
debug:
	go install -ldflags $(LDFLAGS_DEBUG)
	kill -USR1 `pidof autimaat`


## Rebuild the binary and copy it to the remote host.
## Fork the currently running process, if applicable.
##
## Before sending the new binary, we need to delete old one.
## Not doing this, results in scp reporting a "Text file busy" error.
## Overwriting the binary for a running process is apparently not 
## allowed on this host.
release:
	go build -ldflags $(LDFLAGS_RELEASE)
	ssh $(AUTIMAAT_HOST) -p $(AUTIMAAT_PORT) rm $(AUTIMAAT_DIR)/autimaat
	scp -P $(AUTIMAAT_PORT) autimaat $(AUTIMAAT_HOST):$(AUTIMAAT_DIR)
	ssh $(AUTIMAAT_HOST) -p $(AUTIMAAT_PORT) $(AUTIMAAT_DIR)/restart.sh
	go clean