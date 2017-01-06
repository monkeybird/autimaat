// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"os"
	"time"
)

// PayloadHandler defines a function which handles incoming
// server messages.
type PayloadHandler func([]byte)

// ConnectionTimeout defines the deadline for a connection.
const ConnectionTimeout = time.Minute * 10

// Client defines an IRC client for a single network connection.
type Client struct {
	handler PayloadHandler
	conn    net.Conn
	reader  *bufio.Reader
}

// NewClient creates a new client for the given handler.
func NewClient(handler PayloadHandler) *Client {
	return &Client{
		handler: handler,
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
	return c.conn.Close()
}

// File returns the network's file descriptor.
// This call is only valid as long as the connection is actually open.
func (c *Client) File() (*os.File, error) {
	return c.conn.(*net.TCPConn).File()
}

// Run starts the message processing loop and does not return for as long
// as there is an open connection.
func (c *Client) Run() error {
	defer c.Close()

	for {
		line, err := c.read()
		if err != nil {
			return err
		}

		go c.handler(line)
	}
}

// Write writes the given message to the underlying stream.
func (c *Client) Write(p []byte) (int, error) {
	if c.conn == nil || len(p) == 0 {
		return 0, io.EOF
	}

	n, err := c.conn.Write(p)
	if err == nil {
		c.conn.SetDeadline(time.Now().Add(ConnectionTimeout))
	}

	return n, err
}

// Read reads the next message from the connection.
// This call blocks until enough data is available or an error occurs.
func (c *Client) read() ([]byte, error) {
	if c.conn == nil {
		return nil, io.EOF
	}

	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	c.conn.SetDeadline(time.Now().Add(ConnectionTimeout))
	return bytes.TrimSpace(data), nil
}
