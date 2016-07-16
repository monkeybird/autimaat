// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package url scans server messages for URLs. It retrieves
// page titles and returns them to the channel from whence the url came.
package url

import (
	"bytes"
	"fmt"
	"google/youtube"
	"html"
	"io"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/tr"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	// regUrl is used by readUrl to extract web page URLs from incoming
	// PRIVMSG contents.
	regUrl = regexp.MustCompile(`\bhttps?\://[a-zA-Z0-9\-\.]+\.[a-zA-Z]+(\:[0-9]+)?(/\S*)?\b`)

	// regYoutube attempts to look for the video ID part of a youtube url.
	regYoutube = regexp.MustCompile(`[?&]v=([a-zA-Z0-9_-]+)`)

	// youtubeHosts defines a list of known, valid youtube hosts.
	youtubeHosts = []string{
		"youtube.com",
	}

	// These values are used to extract title contents from HTML.
	bOpenTitle1 = []byte("<title>")
	bOpenTitle2 = []byte("<title ")
	bCloseTitle = []byte("</title>")
	bCloseTag   = []byte(">")
)

type module struct {
	youtubeApiKeyFunc func() string
}

// New returns a new module.
func New() mod.Module { return &module{} }

// Load initializes the library and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)
	m.youtubeApiKeyFunc = prof.YoutubeApiKey
}

// Unload cleans up any library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Unbind("PRIVMSG", m.onPrivMsg)
	m.youtubeApiKeyFunc = nil
}

// Help displays help on custom commands.
func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {}

// onPrivMsg checks the given request for any URLs. When found, it returns to
// the channel the title of the page being linked to. This only affects
// resources with content type: text/html.
func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	// Find all URLs in the message body.
	list := regUrl.FindAllString(r.Data, -1)
	if len(list) == 0 {
		return
	}

	// Fetch title data for each of them.
	for _, url := range list {
		go m.fetchTitle(w, r, url)
	}
}

// fetchTitle attempts to retrieve the title element for a given url.
func (m *module) fetchTitle(w irc.ResponseWriter, r *irc.Request, url string) {
	// Ensure the url targets a HTML page. We do this by issueing a HEAD
	// request and checking its content type header.
	resp, err := http.Head(url)
	if err != nil {
		return
	}

	resp.Body.Close()

	ctype := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Index(ctype, "text/html") == -1 &&
		strings.Index(ctype, "text/xhtml") == -1 {
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
	// 16kB is a chunky buffer, but some sites packa a ludicrous amount of
	// crud in their page headers, before getting to the <title> tag.
	var buf [1024 * 16]byte

	// Read the body.
	n, err := io.ReadFull(resp.Body, buf[:])
	resp.Body.Close()

	if err != nil && n <= 0 {
		return // Exit only if no data was read at all.
	}

	body := buf[:n]

	// Extract the title.
	s := bytes.Index(bytes.ToLower(body), bOpenTitle1)
	if s == -1 {
		// title could be something like:
		//
		//    <title xml:lang="en-US">....</title>
		//
		s = bytes.Index(bytes.ToLower(body), bOpenTitle2)
		if s == -1 {
			return
		}

		body = body[s+len(bOpenTitle2):]

		s = bytes.Index(body, bCloseTag)
		if s == -1 {
			return
		}

		body = body[s+1:]
	} else {
		body = body[s+len(bOpenTitle1):]
	}

	e := bytes.Index(bytes.ToLower(body), bCloseTitle)
	if e == -1 {
		e = len(body) - 1
	}

	body = bytes.TrimSpace(body[:e])
	if len(body) == 0 {
		return
	}

	title := html.UnescapeString(string(body))

	// If we are dealing with a youtube link, try to fetch the
	// avideo duration and append it to our response.
	if id, ok := isYoutube(url); ok {
		info, err := youtube.GetVideoInfo(m.youtubeApiKeyFunc(), id)
		if err == nil {
			title += fmt.Sprintf(tr.UrlYoutubeDuration, info.Duration)
		}
	}

	// Show the title to the channel from whence the URL came.
	proto.PrivMsg(w, r.Target, tr.UrlDisplayText, r.SenderName, title)
}

// isYoutube returns a video ID and true if v denotes a youtube video URL.
// Returns false otherwise.
func isYoutube(v string) (string, bool) {
	u, err := url.Parse(v)
	if err != nil {
		return "", false
	}

	if !isYoutubeHost(u.Host) {
		return "", false
	}

	id := strings.TrimSpace(u.Query().Get("v"))
	return id, len(id) > 0
}

// isYoutubeHost returns true if the given value represents a known youtube host.
func isYoutubeHost(v string) bool {
	v = strings.ToLower(v)

	for _, vv := range youtubeHosts {
		if strings.HasSuffix(v, vv) {
			return true
		}
	}

	return false
}
