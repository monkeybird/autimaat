// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package youtube provides a minimal set of bindings for Youtube's
// Data API v3.
//
// This package requires a valid API key to be specified. You can get
// one from your Google account's developer console. See:
//
//    https://console.developers.google.com/apis
//
package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	ErrNoSuchVideo   = errors.New("no such video")
	ErrInvalidAPIKey = errors.New("invalid or missing API key")
	ErrInvalidID     = errors.New("invalid or missing video ID")
)

// VideoInfo defines some detailed properties for a specific
// youtube video.
type VideoInfo struct {
	ID       string        // The video's ID.
	Duration time.Duration // Duration of the video.
}

// GetVideoInfo returns details about a specific video, identified by
// its ID. This is part of the youtube URL. E.g.:
//
//    URL: https://www.youtube.com/watch?v=dQw4w9WgXcQ
//    ID: dQw4w9WgXcQ
//
// This returns nil and an error if the query failed.
func GetVideoInfo(apiKey, id string) (*VideoInfo, error) {
	const videoURL = "https://www.googleapis.com/youtube/v3/videos?id=%s&part=contentDetails&key=%s"

	apiKey = strings.TrimSpace(apiKey)
	apiKey = url.QueryEscape(apiKey)

	if len(apiKey) == 0 {
		return nil, ErrInvalidAPIKey
	}

	id = strings.TrimSpace(id)
	id = url.QueryEscape(id)

	if len(id) == 0 {
		return nil, ErrInvalidID
	}

	var resp videoListResponse
	url := fmt.Sprintf(videoURL, id, apiKey)
	err := fetch(url, &resp)

	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, ErrNoSuchVideo
	}

	item := resp.Items[0]

	return &VideoInfo{
		ID:       item.ID,
		Duration: parseISO8601(item.ContentDetails.Duration),
	}, nil
}

// videoListResponse defines response data for a videoList request.
type videoListResponse struct {
	Items []struct {
		ID             string `json:"id"`
		ContentDetails struct {
			Duration string `json:"duration"` // ISO 8601 timestamp (e.g.: "PT4M13S")
		} `json:"contentDetails"`
	} `json:"items"`
}

// fetch performs an API query and unmarshals the result into the
// given value. Returns an error if something went booboo.
func fetch(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// parseISO8601 parses the given ISO 8601 value into a time.Duration value.
// Example value: "PT4M13S"
//
// ref: https://en.wikipedia.org/wiki/ISO_8601#Durations
func parseISO8601(v string) time.Duration {
	v = strings.ToUpper(v)
	v = strings.TrimSpace(v)

	if len(v) < 3 || v[0] != 'P' {
		return 0
	}

	var err error
	var sum time.Duration
	var inDate bool
	var n int64

	digits := make([]rune, 0, len(v)/2)
	for _, r := range v {
		switch r {
		case 'P':
			inDate = true

		case 'T':
			inDate = false

		case 'Y':
			n, err = strconv.ParseInt(string(digits), 10, 32)
			sum += time.Duration(n) * time.Hour * 8760
			digits = digits[:0]

		case 'M':
			if inDate { // M == month
				n, err = strconv.ParseInt(string(digits), 10, 32)
				sum += time.Duration(n) * time.Hour * 730
				digits = digits[:0]
			} else { // M = minutes
				n, err = strconv.ParseInt(string(digits), 10, 32)
				sum += time.Duration(n) * time.Minute
				digits = digits[:0]
			}

		case 'W':
			n, err = strconv.ParseInt(string(digits), 10, 32)
			sum += time.Duration(n) * time.Hour * 168
			digits = digits[:0]

		case 'D':
			n, err = strconv.ParseInt(string(digits), 10, 32)
			sum += time.Duration(n) * time.Hour * 24
			digits = digits[:0]

		case 'H':
			n, err = strconv.ParseInt(string(digits), 10, 32)
			sum += time.Duration(n) * time.Hour
			digits = digits[:0]

		case 'S':
			n, err = strconv.ParseInt(string(digits), 10, 32)
			sum += time.Duration(n) * time.Second
			digits = digits[:0]

		default:
			if unicode.IsDigit(r) {
				digits = append(digits, r)
			} else {
				return 0
			}
		}

		if err != nil {
			return 0
		}
	}

	return sum
}
