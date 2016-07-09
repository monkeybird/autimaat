// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

import "io"

// ResponseWriter repersents a network stream, used to write
// response data to.
type ResponseWriter interface {
	io.WriteCloser
}
