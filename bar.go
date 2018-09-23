package cmpb

import (
	"bytes"
	"strings"
)

// Bar represents a single progress bar
type Bar struct {
	curr, total int
	lastRender  bool

	p *Progress
}

func newBar(p *Progress, total int) *Bar {
	return &Bar{total: total, p: p}
}

// Update updates the current status of the bar
func (b *Bar) Update(curr int) {
	if b.curr < b.total {
		if curr <= b.total {
			b.curr = curr
		} else {
			b.curr = b.total
		}
		if b.curr == b.total {
			b.lastRender = true
		}
	}
}

// Increment updates the current status of the bar by 1
func (b *Bar) Increment() {
	b.Update(b.curr + 1)
}

func (b *Bar) String() string {
	buf := bytes.Buffer{}
	buf.Grow(b.p.width + 2)
	buf.WriteRune(b.p.lBracket)

	full := b.curr * b.p.width / b.total
	empty := b.p.width - full
	if full > 0 {
		if empty > 0 {
			full--
		}
		buf.WriteString(strings.Repeat(string(b.p.full), full))
		if empty > 0 {
			buf.WriteRune(b.p.curr)
		}
	}
	buf.WriteString(strings.Repeat(string(b.p.empty), empty))

	buf.WriteRune(b.p.rBracket)
	return buf.String()
}
