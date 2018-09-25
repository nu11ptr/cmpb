package cmpb

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/nu11ptr/cmpb/strutil"
)

// Bar represents a single progress bar
type Bar struct {
	key, action       string
	curr, total       int
	lastRender        bool
	start             time.Time
	preBarF, postBarF func(int, int, time.Time) string

	p   *Progress
	mut sync.Mutex
}

func newBar(key string, total int, p *Progress) *Bar {
	return &Bar{
		key: key, action: "", total: total, start: time.Now(), preBarF: calcDur(), postBarF: calcPct, p: p,
	}
}

func calcDur() func(int, int, time.Time) string {
	var final time.Time

	return func(curr, total int, start time.Time) string {
		last := time.Now()

		// Record the time a single time on the last update
		if curr == total {
			if final.IsZero() {
				final = last
			}
			last = final
		}
		dur := last.Sub(start)
		dur = dur.Round(time.Second)

		var buf bytes.Buffer
		if dur >= time.Hour {
			h := dur / time.Hour
			dur -= h * time.Hour
			buf.WriteString(fmt.Sprintf("%dh", h))
		}
		if dur >= time.Minute {
			m := dur / time.Minute
			dur -= m * time.Minute
			buf.WriteString(fmt.Sprintf("%dm", m))
		}
		if dur >= time.Second {
			s := dur / time.Second
			dur -= s * time.Second
			buf.WriteString(fmt.Sprintf("%ds", s))
		}
		return color.HiMagentaString(buf.String())
	}
}

func calcPct(curr, total int, start time.Time) string {
	pct := (curr * 100) / total
	return fmt.Sprintf(color.HiMagentaString("%d%%"), pct)
}

func (b *Bar) update(curr int) {
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

// SetAction sets the displayed current action
func (b *Bar) SetAction(action string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.action = action
}

// Update updates the current status of the bar
func (b *Bar) Update(curr int) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.update(curr)
}

// Increment updates the current status of the bar by 1
func (b *Bar) Increment() {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.update(b.curr + 1)
}

func (b *Bar) isLastRender() bool {
	b.mut.Lock()
	defer b.mut.Unlock()

	lr := b.lastRender
	if lr {
		b.lastRender = false
	}
	return lr
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
	b.mut.Lock()
	defer b.mut.Unlock()

	buf := new(bytes.Buffer)
	param := &b.p.param
	w := param.PrePad + param.KeyWidth + len(param.KeyDiv) + param.ActionWidth + param.PreBarWidth +
		param.BarWidth + param.PostBarWidth + 4
	buf.Grow(w)

	buf.WriteString(strings.Repeat(" ", param.PrePad))
	buf.WriteString(strutil.ResizeR(b.key, param.Post, param.KeyWidth))
	buf.WriteString(param.KeyDiv)
	buf.WriteRune(' ')
	buf.WriteString(strutil.ResizeR(b.action, param.Post, param.ActionWidth))
	buf.WriteRune(' ')
	buf.WriteString(strutil.ResizeL(b.preBarF(b.curr, b.total, b.start), param.Post, param.PreBarWidth))
	buf.WriteRune(' ')
	b.makeBar(param, buf)
	buf.WriteRune(' ')
	buf.WriteString(strutil.ResizeL(b.postBarF(b.curr, b.total, b.start), param.Post, param.PostBarWidth))

	return buf.String()
}
