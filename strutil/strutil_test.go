package strutil_test

import (
	"testing"

	"github.com/nu11ptr/cmpb/strutil"
)

func TestLen(t *testing.T) {
	tests := []struct {
		name, input string
		output      int
	}{
		{"Empty", "", len("")},
		{"Plain", "abc abc", len("abc abc")},
		{"FalsePositive", "abc\x1babc", len("abc abc")},
		{"SetAndResetColor", "\x1b[36mabc abc\x1b[0m", len("abc abc")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := strutil.Len(test.input)
			if l != test.output {
				t.Error("want", test.output, "got", l)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name, input, output string
		l                   int
		trunc               bool
	}{
		{"TruncNotNeeded", "abc", "abc", 3, false},
		{"Basic", "abcabc", "abc", 3, true},
		{"HasEscapes", "abc\x1b[36mabc", "abc\x1b[36mab", 5, true},
		{"FalsePositive", "abc\x1babc", "abc\x1ba", 5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s2, trunc := strutil.Truncate(test.input, test.l)
			if s2 != test.output {
				t.Error("want", test.output, "got", s2)
			}
			if trunc != test.trunc {
				t.Error("want", test.trunc, "got", trunc)
			}
		})
	}
}
