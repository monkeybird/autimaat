// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package mod defines the interface for command modules. There are
// essentially plugins which provide specific functionality for the bot.
package mod

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"os"
)

// Module defines the interface for a single module.
type Module interface {
	// Load initializes the module and loads any internal resources
	// which may be required.
	Load(irc.ProtocolBinder, irc.Profile)

	// Unload cleans the module up and unloads any internal resources.
	Unload(irc.ProtocolBinder, irc.Profile)

	// Help is called whenever the module is to provide help on any
	// custom commands it may implement.
	Help(irc.ResponseWriter, *cmd.Request)
}

// Load loads the, optionally compressed, contents of the given
// file and unmarshals it into the specified value.
func Load(file string, v interface{}, compressed bool) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	var r io.Reader
	if compressed {
		gz, err := gzip.NewReader(fd)
		if err != nil {
			return err
		}

		defer gz.Close()
		r = gz
	} else {
		r = fd
	}

	return json.NewDecoder(r).Decode(v)
}

// Save writes the marshaled version of v to the given file.
// It is optionally gzip compressed.
func Save(file string, v interface{}, compressed bool) error {
	if compressed {
		fd, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}

		defer fd.Close()

		gz := gzip.NewWriter(fd)
		defer gz.Close()

		return json.NewEncoder(gz).Encode(v)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0600)
}
