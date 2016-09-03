// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"regexp"
)

var (
	RegAny     = regexp.MustCompile(`^.*$`)
	RegInt     = regexp.MustCompile(`^[+-]?\d+$`)
	RegUint    = regexp.MustCompile(`^[+]?\d+$`)
	RegFloat   = regexp.MustCompile(`^[+-]?\d+(\.\d+([eE][+-]?\d+)?)?$`)
	RegBool    = regexp.MustCompile(`^(1|0|t(rue)?|f(alse)?|y(es)?|no?|on|off)$`)
	RegChannel = regexp.MustCompile(`^[#&+!][^ ,:]{1,50}$`)
	RegMode    = regexp.MustCompile(`^[+-][obveI]$`)
	RegUrl     = regexp.MustCompile(`^https?\://[a-zA-Z0-9\-\.]+\.[a-zA-Z]+(\:[0-9]+)?(/\S*)?$`)
)
