package strutil_test

import (
	"testing"
	"time"

	"github.com/fatih/color"
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

func TestResize(t *testing.T) {
	post := color.HiCyanString("...")
	tests := []struct {
		name, input, output, post string
		l                         int
		f                         func(string, string, int) string
	}{
		{"JustRight", color.HiCyanString("abc"), color.HiCyanString("abc"), post, 3, strutil.ResizeL},
		{"TooShortL", color.HiCyanString("abc"), "   " + color.HiCyanString("abc"), post, 6, strutil.ResizeL},
		{"TooShortR", color.HiCyanString("abc"), color.HiCyanString("abc") + "   ", post, 6, strutil.ResizeR},
		{"TooLong", color.HiCyanString("abcabcabc"), color.HiCyanString("abc") + post, post, 6, strutil.ResizeL},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.f(test.input, test.post, test.l)
			if output != test.output {
				t.Error("want", test.output, "got", output)
			}
		})
	}
}

func TestFmtDuration(t *testing.T) {
	tests := []struct {
		name, output string
		input        time.Duration
	}{
		{"Zero", "0s", time.Duration(0)},
		{"MinSec", "45m 23s", time.Duration((45 * time.Minute) + (23 * time.Second))},
		{"HourMinSec", "67h 45m 23s",
			time.Duration((67 * time.Hour) + (45 * time.Minute) + (23 * time.Second))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := strutil.FmtDuration(test.input)
			if output != test.output {
				t.Error("want", test.output, "got", output)
			}
		})
	}
}
