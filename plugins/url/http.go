// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins/url/youtube"
)

var (
	// regUrl is used by readUrl to extract web page URLs from incoming
	// PRIVMSG contents.
	regUrl = regexp.MustCompile(`\bhttps?\://[a-zA-Z0-9\-\.]+\.[a-zA-Z]+(\:[0-9]+)?(/\S*)?\b`)

	// These values are used to extract title contents from HTML.
	bOpenTitle1 = []byte("<title>")
	bOpenTitle2 = []byte("<title ")
	bCloseTitle = []byte("</title>")
	bCloseTag   = []byte(">")
)

// fetchTitle attempts to retrieve the title element for a given url.
func fetchTitle(w irc.ResponseWriter, r *irc.Request, url, apiKey string) {
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
	// video duration and append it to our response.
	if id := isYoutube(url); len(id) > 0 {
		info, err := youtube.GetVideoInfo(apiKey, id)
		if err == nil {
			title += fmt.Sprintf(TextYoutubeDuration, info.Duration)
		}
	}

	// If the title matches one of the titles that we want to ignore,
	// do not show it.
	if Ignore[title] {
		return
	}

	// Show the title to the channel from whence the URL came.
	proto.PrivMsg(w, r.Target, TextDisplay, r.SenderName, title)
}

// isYoutube returns a video ID and true if v denotes a recognized youtube
// video URL. Returns an empty string otherwise.
func isYoutube(v string) string {
	u, err := url.Parse(v)
	if err != nil {
		return ""
	}

	if strings.EqualFold(u.Host, "youtube.com") ||
		strings.EqualFold(u.Host, "www.youtube.com") {
		id := strings.TrimSpace(u.Query().Get("v"))
		return id
	}

	if strings.EqualFold(u.Host, "youtu.be") ||
		strings.EqualFold(u.Host, "www.youtu.be") {
		id := u.RequestURI()
		if strings.HasPrefix(id, "/") {
			id = id[1:]
		}
		return id
	}

	return ""
}
