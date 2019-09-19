package misc_test

import (
	"testing"
	"unicode/utf8"

	"github.com/setlog/trivrost/pkg/misc"
)

func TestMustGetRandomHexString(t *testing.T) {
	for i := 0; i < 10; i++ {
		for l := 0; l <= 16; l++ {
			s := misc.MustGetRandomHexString(l)
			length := utf8.RuneCountInString(s)
			if length != 2*l {
				t.Errorf("String %s length was %d. Expected %d.", s, length, 2*l)
			}
			for p, r := range s {
				if r >= '0' && r <= '9' {
				} else if r >= 'a' && r <= 'f' {
				} else if r >= 'A' && r <= 'F' {
				} else {
					t.Errorf("Non-hex symbol %c at index %d for test run %d with bytecount %d.", r, p, i, l)
				}
			}
		}
	}
}
