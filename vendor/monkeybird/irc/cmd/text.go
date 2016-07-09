// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import "strings"

// split splits the given string into individual parameters.
// This takes quoted strings into account.
func split(data string) (string, []string) {
	var quoted bool
	var start, pos int

	set := make([]string, 0, 5)

	for pos = 0; pos < len(data); pos++ {
		switch data[pos] {
		case '"':
			quoted = !quoted

		case ' ', '\t':
			if quoted {
				break
			}

			str := strings.TrimSpace(data[start:pos])

			if len(str) > 0 {
				str = strings.Replace(str, "\"", "", -1)
				set = append(set, str)
			}

			start = pos + 1
		}
	}

	if start < pos {
		str := strings.TrimSpace(data[start:pos])
		if len(str) > 0 {
			str = strings.Replace(str, "\"", "", -1)
			set = append(set, str)
		}
	}

	if len(set) == 0 {
		return "", nil
	}

	return set[0], set[1:]
}
