// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package youtube

import (
	"testing"
	"time"
)

// Put your youtube API key here.
// see: https://console.developers.google.com/apis
const ApiKey = ""

func TestParseISO8601(t *testing.T) {
	testISO8601(t, "", 0)
	testISO8601(t, "PY", 0)
	testISO8601(t, "PT", 0)
	testISO8601(t, "P1Y", time.Hour*8760)
	testISO8601(t, "P1H", time.Hour)
	testISO8601(t, "PT1H", time.Hour)
	testISO8601(t, "PTz1H", 0)
	testISO8601(t, "P 1H", 0)
	testISO8601(t, "P2MT2M33S", 2*time.Hour*730+2*time.Minute+33*time.Second)
}

func testISO8601(t *testing.T, in string, want time.Duration) {
	have := parseISO8601(in)
	if want != have {
		t.Fatalf("ISO 8601 mismatch for %q; want: %s\nhave %s", in, want, have)
	}
}

func TestGetVideoInfo(t *testing.T) {
	if len(ApiKey) == 0 {
		return
	}

	const videoID = "dQw4w9WgXcQ"

	info, err := GetVideoInfo(ApiKey, videoID)
	if err != nil {
		t.Fatal(err)
	}

	if info.ID != videoID {
		t.Fatalf("Video ID mismatch; want %q\nhave: %q", videoID, info.ID)
	}
}
