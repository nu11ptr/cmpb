package cmpb

import (
	"bytes"
	"strings"

	"github.com/nu11ptr/cmpb/strutil"
)

// Bar represents a single progress bar
type Bar struct {
	key         string
	curr, total int
	lastRender  bool

	p *Progress
}

func newBar(key string, total int, p *Progress) *Bar {
	return &Bar{key: key, total: total, p: p}
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

func (b *Bar) makeBar(param *Param, buf *bytes.Buffer) {
	buf.WriteString(param.LBracket)

	full := b.curr * (param.BarWidth - 2) / b.total
	empty := (param.BarWidth - 2) - full
	if full > 0 {
		if empty > 0 {
			full--
		}
		buf.WriteString(strings.Repeat(string(param.Full), full))
		if empty > 0 {
			buf.WriteString(param.Curr)
		}
	}
	buf.WriteString(strings.Repeat(string(param.Empty), empty))

	buf.WriteString(param.RBracket)
}

func (b *Bar) String() string {
	buf := new(bytes.Buffer)
	param := &b.p.param
	w := param.PrePad + param.KeyWidth + len(param.KeyDiv) + param.ActionWidth + param.PreBarWidth +
		param.BarWidth + param.PostBarWidth + 4
	buf.Grow(w)

	buf.WriteString(strings.Repeat(" ", param.PrePad))
	buf.WriteString(strutil.Resize(b.key, param.Post, param.KeyWidth))
	buf.WriteString(param.KeyDiv)
	buf.WriteRune(' ')
	buf.WriteString(strutil.Resize("downloading...", param.Post, param.ActionWidth))
	buf.WriteRune(' ')
	buf.WriteString(strutil.Resize("00h00m00s", param.Post, param.PreBarWidth))
	buf.WriteRune(' ')
	b.makeBar(param, buf)
	buf.WriteRune(' ')
	buf.WriteString(strutil.Resize("100%", param.Post, param.PostBarWidth))

	return buf.String()
}
