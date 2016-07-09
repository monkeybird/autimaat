// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// HandleFunc defines a function which handles incoming server messages.
type HandleFunc func(*Request)

// Timeout deadline for connection.
const Timeout = time.Minute * 10

// Client defines an IRC client for a single network connection.
type Client struct {
	profile  Profile
	handle   HandleFunc
	conn     net.Conn
	reader   *bufio.Reader
	quitOnce sync.Once
}

// NewClient creates a new client for the given profile.
func NewClient(profile Profile, handle HandleFunc) *Client {
	return &Client{
		profile: profile,
		handle:  handle,
	}
}

// Open creates a new client connection to the given address with the format:
// <host>:<port>.
//
// If the tls config is not nil, it will be used to upgrade the connection
// to a TLS connection.
func (c *Client) Open(address string, cfg *tls.Config) error {
	var err error

	c.conn, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}

	if cfg != nil {
		c.reader = bufio.NewReader(tls.Client(c.conn, cfg))
	} else {
		c.reader = bufio.NewReader(c.conn)
	}

	return nil
}

// OpenFd opens a new client from the given file descriptor.
// If the tls config is not nil, it will be used to upgrade the connection
// to a TLS connection.
func (c *Client) OpenFd(file *os.File, cfg *tls.Config) error {
	var err error

	c.conn, err = net.FileConn(file)
	if err != nil {
		return err
	}

	if cfg != nil {
		c.reader = bufio.NewReader(tls.Client(c.conn, cfg))
	} else {
		c.reader = bufio.NewReader(c.conn)
	}

	return nil
}

// Close closes the connection.
func (c *Client) Close() error {
	c.quitOnce.Do(func() {
		c.conn.Close()
	})
	return nil
}

// File returns the network's file descriptor.
// This call is only valid as long as the connection is actually open.
func (c *Client) File() (*os.File, error) {
	a, b := c.conn.(*net.TCPConn).File()
	return a, b
}

// Run starts the message processing loop and does not return for as long
// as there is an open connection.
func (c *Client) Run() error {
	var line []byte
	var err error

	defer c.Close()
	prof := c.profile

	for {
		line, err = c.read()
		if err != nil {
			return err
		}

		r := parseRequest(line)
		if r != nil {
			// If Target points to the bot's own name, this is from a user, not
			// a channel. Change the Target to the sender's name, so any replies
			// we create, end up at the right destination.
			if prof.IsNick(r.Target) {
				r.Target = r.SenderName
			}

			c.handle(r)
		}
	}

	return nil
}

// Write writes the given message to the underlying stream.
func (c *Client) Write(p []byte) (int, error) {
	if c.conn == nil || len(p) == 0 {
		return 0, io.EOF
	}

	n, err := c.conn.Write(p)
	c.conn.SetDeadline(time.Now().Add(Timeout))
	return n, err
}

// Read reads the next message from the connection.
// This call blocks until enough data is available or an error occurs.
func (c *Client) read() ([]byte, error) {
	conn := c.conn
	rdr := c.reader

	if conn == nil {
		return nil, io.EOF
	}

	data, err := rdr.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(Timeout))
	return bytes.TrimSpace(data), nil
}
