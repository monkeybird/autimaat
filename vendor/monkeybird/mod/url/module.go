// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package url scans server messages for URLs. It retrieves
// page titles and returns them to the channel from whence the url came.
package url

import (
	"bytes"
	"html"
	"io"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/tr"
	"net/http"
	"regexp"
	"strings"
)

var (
	// regUrl is used by readUrl to extract web page URLs from incoming
	// PRIVMSG contents.
	regUrl = regexp.MustCompile(`\bhttps?\://[a-zA-Z0-9\-\.]+\.[a-zA-Z]+(\:[0-9]+)?(/\S*)?\b`)

	// These values are used to extract title contents from HTML.
	bOpenTitle  = []byte("<title>")
	bCloseTitle = []byte("</title>")
)

type module struct{}

func New() mod.Module {
	return &module{}
}

// Load initializes the library and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", onPrivMsg)
}

// Unload cleans up any library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Unbind("PRIVMSG", onPrivMsg)
}

// Help displays help on custom commands.
func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {}

// onPrivMsg checks the given request for any URLs. When found, it returns to
// the channel the title of the page being linked to. This only affects
// resources with content type: text/html.
func onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	if !r.FromChannel() {
		return
	}

	// Find all URLs in the message body.
	list := regUrl.FindAllString(r.Data, -1)
	if len(list) == 0 {
		return
	}

	// Fetch title data for each of them.
	for _, url := range list {
		go fetchTitle(w, r, url)
	}
}

// fetchTitle attempts to retrieve the title element for a given url.
func fetchTitle(w irc.ResponseWriter, r *irc.Request, url string) {
	// Ensure the url targets a HTML page. We do this by issueing a HEAD
	// request and checking its content type header.
	resp, err := http.Head(url)
	if err != nil {
		return
	}

	resp.Body.Close()

	ctype := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Index(ctype, "text/html") == -1 {
		return
	}

	// We have an HTML document -- Fetch its contents.
	resp, err = http.Get(url)
	if err != nil {
		return
	}

	// buf defines the maximum amount of data we will be reading from a page,
	// before stopping our search for the <title> tag.
	//
	// 16kB is pretty large, but some larger sites pack ludicrous amounts
	// of crud in their page headers, before getting to the title.
	var buf [1024 * 16]byte

	// Read the body.
	n, err := io.ReadFull(resp.Body, buf[:])
	resp.Body.Close()
	if err != nil {
		return
	}

	body := buf[:n]

	// Extract the title.
	s := bytes.Index(bytes.ToLower(body), bOpenTitle)
	if s == -1 {
		return
	}

	body = body[s+7:]

	e := bytes.Index(bytes.ToLower(body), bCloseTitle)
	if e == -1 {
		e = len(body) - 1
	}

	body = bytes.TrimSpace(body[:e])
	if len(body) == 0 {
		return
	}

	// Show the title to the channel from whence the URL came.
	proto.PrivMsg(w, r.Target, tr.UrlDisplayText,
		r.SenderName, html.UnescapeString(string(body)))
}
