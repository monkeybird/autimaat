// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package util defines a few commonly used utility functions.
package util

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Action returns the given, formatted message as a user action.
// E.g.: /me <something something>
func Action(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x01ACTION %s\x01", fmt.Sprintf(f, argv...))
}

// Bold returns the given value as bold text.
func Bold(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x02%s\x02", fmt.Sprintf(f, argv...))
}

// Italic returns the given value as italicized text.
func Italic(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x1d%s\x1d", fmt.Sprintf(f, argv...))
}

// Underline returns the given value as underlined text.
func Underline(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x1f%s\x1f", fmt.Sprintf(f, argv...))
}

// ReadFile loads the, optionally compressed, contents of the given
// file and unmarshals it into the specified value v.
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
