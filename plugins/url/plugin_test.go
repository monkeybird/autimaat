// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestYoutube(t *testing.T) {
	testYoutube(t, "boops", "")
	testYoutube(t, "youtube.com", "")
	testYoutube(t, "http://youtube.com", "")
	testYoutube(t, "https://www.youtube.com", "")
	testYoutube(t, "https://www.not-youtube.com", "")
	testYoutube(t, "https://www.not-youtube.com?v=boops", "")
	testYoutube(t, "https://youtube.com?v=HKNXXpqareI", "HKNXXpqareI")
	testYoutube(t, "https://www.youtube.com?v=HKNXXpqareI", "HKNXXpqareI")
	testYoutube(t, "http://youtube.com?v=HKNXXpqareI", "HKNXXpqareI")
	testYoutube(t, "http://www.youtube.com?v=HKNXXpqareI", "HKNXXpqareI")
	testYoutube(t, "https://www.youtube.com/watch?v=gc1VA_W3M-Y&index=2&list=PLJwv6sN_mnF0QsOTcKlFDeyzwXMM0MWru", "gc1VA_W3M-Y")
}

func testYoutube(t *testing.T, in, want string) {
	have := isYoutube(in)

	if want != have {
		t.Fatalf("id mismatch for %q;\nwant: %q\nhave: %q",
			in, want, have)
	}
}

func TestTitle(t *testing.T) {
	testTitle(t, false, "https://www.youtube.com/watch?v=BDB4ZF8jX9Q", "Betelgeuse Supernova and Its Impact On Earth - Documentary - YouTube")
	testTitle(t, false, "https://open.spotify.com/track/3wdLmIe8zlifCoOhb3N4nK", "Spotify Web Player - Last Run - Tokyo Rose")
	testTitle(t, false, "https://open.spotify.com/album/2dVQfNoIYa32dyP6PK7zbw", "Spotify Web Player - Space Tapes & Vice - Bourgeoisie")
	testTitle(t, true, "I do not exist.", "")
	testTitle(t, true, "https://i.imgur.com/VZrAPSv.gif", "")
	testTitle(t, false, "https://imgur.com/gallery/rRqwy", "My current relationship with Imgur. - Album on Imgur")
	testTitle(t, false, "https://open.spotify.com/track/75mx4MRQt4l7Gs49JSc6QV", "Spotify Web Player - Through the Barricades - Various Artists")
	testTitle(t, false, "http://rtl.nl", "RTL XL")
	testTitle(t, false, "https://en.wikipedia.org/wiki/Neolithic_Subpluvial", "Neolithic Subpluvial - Wikipedia")
}

func testTitle(t *testing.T, expectError bool, url, want string) {
	have, err := fetchTitle(url, "")
	if !expectError && err != nil {
		t.Fatalf("unexpected error for %q:\nerror: %v", url, err)
	}

	if expectError && err == nil {
		t.Fatalf("expected error for %q:\ngot title: %q", url, have)
	}

	if !strings.EqualFold(have, want) {
		t.Fatalf("title mismatch for %q:\nhave: %q\nwant: %q", url, have, want)
	}
}
