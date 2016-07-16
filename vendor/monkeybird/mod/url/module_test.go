// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

import (
	"monkeybird/irc"
	"os"
	"testing"
)

var (
	bindings irc.BindingList
	profile  irc.Profile
	tm       *module
)

func TestMain(m *testing.M) {
	profile = irc.NewProfile(".")

	tm = New().(*module)
	tm.Load(&bindings, profile)
	ret := m.Run()
	tm.Unload(&bindings, profile)

	os.Exit(ret)
}

func TestYoutube(t *testing.T) {
	testYoutube(t, "boops", "", false)
	testYoutube(t, "youtube.com", "", false)
	testYoutube(t, "http://youtube.com", "", false)
	testYoutube(t, "https://www.youtube.com", "", false)
	testYoutube(t, "https://www.not-youtube.com", "", false)
	testYoutube(t, "https://www.not-youtube.com?v=boops", "boops", true)
	testYoutube(t, "https://youtube.com?v=HKNXXpqareI", "HKNXXpqareI", true)
	testYoutube(t, "https://www.youtube.com?v=HKNXXpqareI", "HKNXXpqareI", true)
	testYoutube(t, "http://youtube.com?v=HKNXXpqareI", "HKNXXpqareI", true)
	testYoutube(t, "http://www.youtube.com?v=HKNXXpqareI", "HKNXXpqareI", true)
	testYoutube(t, "https://www.youtube.com/watch?v=gc1VA_W3M-Y&index=2&list=PLJwv6sN_mnF0QsOTcKlFDeyzwXMM0MWru", "gc1VA_W3M-Y", true)
}

func testYoutube(t *testing.T, in, wantA string, wantB bool) {
	haveA, haveB := isYoutube(in)

	if wantA != haveA {
		t.Fatalf("Youtube id mismatch for %q;\nwant: %q\nhave: %q",
			in, wantA, haveA)
	}

	if wantB != haveB {
		t.Fatalf("Youtube succes mismatch for %q;\nwant: %v\nhave: %v",
			in, wantB, haveB)
	}
}
