// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package plugins defines and maintains supported bot plugins.
package plugins

import (
	"log"

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
