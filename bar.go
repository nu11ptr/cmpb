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
	param := &b.p.param
	buf.Grow(param.Width + 2)
	buf.WriteRune(param.LBracket)

	full := b.curr * param.Width / b.total
	empty := param.Width - full
	if full > 0 {
		if empty > 0 {
			full--
		}
		buf.WriteString(strings.Repeat(string(param.Full), full))
		if empty > 0 {
			buf.WriteRune(param.Curr)
		}
	}
	buf.WriteString(strings.Repeat(string(param.Empty), empty))

	buf.WriteRune(param.RBracket)
	return buf.String()
}
