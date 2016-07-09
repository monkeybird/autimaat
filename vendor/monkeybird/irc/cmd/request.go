// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import "monkeybird/irc"

// Request extends an irc.Request with command context like parsed parameters.
type Request struct {
	*irc.Request
	Params []Param
}

// Len returns the number of parameters passed with this request.
func (r *Request) Len() int            { return len(r.Params) }
func (r *Request) String(n int) string { return r.Params[n].String() }
func (r *Request) Int(n int) int64     { return r.Params[n].Int() }
func (r *Request) Uint(n int) uint64   { return r.Params[n].Uint() }
func (r *Request) Float(n int) float64 { return r.Params[n].Float() }
func (r *Request) Bool(n int) bool     { return r.Params[n].Bool() }
