## autimaat

This program is an IRC bot, specifically written for a private IRC channel.


## Install

    go get github.com/monkeybird/autimaat


## Usage

First, create a new profile directory and configuration file:

	$ autimaat -new /path/to/profile

Edit the newly created configuration file to your liking:

	$ nano /path/to/profile/profile.cfg

Relaunch the bot to use the new profile:

	$ autimaat /path/to/profile

In order to have the bot automatically re-launch after shutdown, an external
supervisor like systemd is required. The bot will create a PID file at
`/path/to/profile/app.pid`, in case the supervisor requires it.

When dealing with systemd, the bot may have to be forked at least once,
after it has been launched. Otherwise, systemd will keep killing it and
re-launching it in a never ending loop. Forking the bot is done through
the following command:

	$ kill -s USR1 `pidof autimaat`

This tells the bot to fork itself, while passing along any existing connections.
The old process then shuts itself down. This mechanism allows the bot to be binary-
patched, without downtime.


## license

Unless otherwise noted, the contents of this project are subject to a 1-clause BSD
license. Its contents can be found in the enclosed LICENSE file.
