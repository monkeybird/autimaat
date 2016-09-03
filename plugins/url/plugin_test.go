// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

import (
	"os"
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
