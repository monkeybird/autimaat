// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

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

	ErrNotHTML = errors.New("url is not an HTML resource")
	ErrNoTitle = errors.New("url contains no title")
)

// fetchTitle attempts to retrieve the title element for a given url.
// Returns an error if something went wrong, or a empty string if the
// title is meant to be filtered out.
func fetchTitle(url, apiKey string) (string, error) {
	var client http.Client

	err := fetchHead(&client, url)
	if err != nil {
		return "", err
	}

	body, err := fetchBody(&client, url)
	if err != nil || len(body) == 0 {
		return "", err
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
		return "", ErrNoTitle
	}

	return title, nil
}

// Ensure the url targets a HTML page. We do this by issueing a HEAD
// request and checking its content type header.
func fetchHead(client *http.Client, url string) error {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	// Set a specific user agent. Some sites, like Spotify will not
	// function properly if no supported user-agent is defined.
	req.Header.Set("User-Agent", TextUserAgent)

	// Ensure the url targets a HTML page. We do this by issueing a HEAD
	// request and checking its content type header.
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	ctype := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Index(ctype, "text/html") == -1 &&
		strings.Index(ctype, "text/xhtml") == -1 {
		return ErrNotHTML
	}

	return nil
}

// fetchBody returns the body of the given url. Or at least a subset of it.
func fetchBody(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set a specific user agent. Some sites, like Spotify will not
	// function properly if no supported user-agent is defined.
	req.Header.Set("User-Agent", TextUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
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

	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	if (err == nil || err == io.ErrUnexpectedEOF) && n <= 0 {
		return nil, ErrNoTitle
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
			return nil, ErrNoTitle
		}

		body = body[s+len(bOpenTitle2):]

		s = bytes.Index(body, bCloseTag)
		if s == -1 {
			return nil, ErrNoTitle
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
		return nil, ErrNoTitle
	}

	return body, nil
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
