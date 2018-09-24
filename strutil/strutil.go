package strutil

import "bytes"

type strState int

const (
	notInEscape strState = iota
	inEscape
	// startingEscape means we have seen the ESC char but are not yet in a full escape sequence
	startingEscape

	esc      = '\x1b'
	lBracket = '['
)

// Len computes the length of a string, but unlike the builtin len, it ignores ANSI escape codes
func Len(s string) (count int) {
	state := notInEscape

	for _, c := range s {
		switch state {
		case notInEscape:
			if c == esc {
				state = startingEscape
			} else {
				count++
			}
		case inEscape:
			if c >= '@' && c <= '~' {
				state = notInEscape
			}
		case startingEscape:
			if c == lBracket {
				state = inEscape
			} else {
				state = notInEscape
				// We increment count because this escape was a false positive and wasn't counted earlier
				// Additionally, this 2nd char (that wasn't lBracket) was also not counted and should be
				count += 2
			}
		}
	}
	return
}

// Truncate truncates the string s to the given length (ignoring ANSI sequences) and returns the
// new string. It also returns a boolean based on whether it actually needed to truncate or not
func Truncate(s string, l int) (string, bool) {
	if Len(s) <= l {
		return s, false
	}
	state := notInEscape
	count := 0
	buf := bytes.Buffer{}
	buf.Grow(len(s)) // The biggest it could possibly get is how big it is now (using builtin len)

	for _, c := range s {
		switch state {
		case notInEscape:
			if c == esc {
				state = startingEscape
			} else if count < l {
				buf.WriteRune(c)
				count++
			}
		case inEscape:
			if c >= '@' && c <= '~' {
				state = notInEscape
			}
			buf.WriteRune(c)
		case startingEscape:
			if c == lBracket {
				state = inEscape
				// We can't write out esc until we get here since we don't know if false positive
				buf.WriteRune(esc)
				buf.WriteRune(c)
			} else {
				state = notInEscape
				// This was a false positive and wasn't counted earlier so we now write out these chars
				if count < l {
					buf.WriteRune(esc)
					count++
				}
				if count < l {
					buf.WriteRune(c)
					count++
				}
			}
		}
	}
	return buf.String(), true
}
