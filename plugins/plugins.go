// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package plugins defines and maintains supported bot plugins.
// Additionally, it provides some rudimentary utility functions
// for plugins to use.
package plugins

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/monkeybird/autimaat/irc"
)

// Plugin defines the interface for a single plugin.
type Plugin interface {
	// Load initializes the module and loads any internal resources
	// which may be required.
	Load(irc.Profile) error

	// Unload cleans the module up and unloads any internal resources.
	Unload(irc.Profile) error

	// Dispatch sends the given, incoming IRC message to the plugin for
	// processing as it sees fit.
	Dispatch(irc.ResponseWriter, *irc.Request)
}

// List of registered plugins. This is to be filled during
// proigram initialization and is considered read-only from then on.
var plugins []Plugin

// Register registers the given plugin. This is meant to be called during
// program initialization, by imported plugin packages.
func Register(p Plugin) { plugins = append(plugins, p) }

// Load initializes all plugins.
func Load(prof irc.Profile) {
	for _, p := range plugins {
		log.Printf("[plugins] Loading: %T", p)

		err := p.Load(prof)
		if err != nil {
			log.Printf("[%T] %v", p, err)
		}
	}
}

// Unload unloads all plugins.
func Unload(prof irc.Profile) {
	for _, p := range plugins {
		log.Printf("[plugins] Unloading: %T", p)

		err := p.Unload(prof)
		if err != nil {
			log.Printf("[%T] %v", p, err)
		}
	}
}

// Dispatch sends the given, incoming IRC message to all plugins.
func Dispatch(w irc.ResponseWriter, r *irc.Request) {
	for _, p := range plugins {
		go p.Dispatch(w, r)
	}
}

// ReadFile loads the, optionally compressed, contents of the given
// file and unmarshals it into the specified value.
//
// This is a utility function which can be used by plugins to load
// custom configuration info or data from disk.
func ReadFile(file string, v interface{}, compressed bool) error {
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

// WriteFile writes the marshaled version of v to the given file.
// It is optionally gzip compressed.
//
// This is a utility function which can be used by plugins to save
// custom configuration info or data to disk.
func WriteFile(file string, v interface{}, compressed bool) error {
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
