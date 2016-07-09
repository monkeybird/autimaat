#!/usr/bin/env sh

## This script rebuilds the bot and replaces any running instance with the
## new binary.
##
## This script is mostly used by me during testing and debugging of new commands.
## It allows me to patch the bot, without having to reconnect to the network.

## rebuild the bot binary and make sure its revision number is incremented.
go install -ldflags "-X main.AppVersionRevision=`date -u +%s`"

## Ignore the rest if this failed for whatever reason.
if [ ! $? -eq 0 ]; then
	exit $?;
fi

## Tell the running bot process it is time to fork itself.
kill -USR1 `pidof autimaat`