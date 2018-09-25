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
	key, msg            string
	curr, total         int
	lastRender, stopped bool
	start               time.Time
	preBarF, postBarF   func(int, int, time.Time, bool) string

	p   *Progress
	mut sync.Mutex
}

func newBar(key string, total int, p *Progress) *Bar {
	return &Bar{
		key: key, msg: "", total: total, start: time.Now(), preBarF: calcDur(), postBarF: calcPct, p: p,
	}
}

func calcDur() func(int, int, time.Time, bool) string {
	var final time.Time

	return func(curr, total int, start time.Time, stopped bool) string {
		last := time.Now()

		// Record the time a single time on the last update so it doesn't keep updating after bar is done
		if curr == total || stopped {
			if final.IsZero() {
				final = last
			}
			last = final
		}
		dur := last.Sub(start)
		return color.HiMagentaString(strutil.FmtDuration(dur))
	}
}

func calcPct(curr, total int, start time.Time, stopped bool) string {
	pct := (curr * 100) / total
	return fmt.Sprintf(color.HiMagentaString("%d%%"), pct)
}

func (b *Bar) update(curr int) {
	if b.curr < b.total && !b.stopped {
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

// Stop stops the updating of the bar and sets a final msg
func (b *Bar) Stop(msg string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.stopped = true
	b.lastRender = true
	b.msg = msg
}

// SetMessage sets the displayed current message
func (b *Bar) SetMessage(msg string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.msg = msg
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
	w := param.PrePad + param.KeyWidth + len(param.KeyDiv) + param.MsgWidth + param.PreBarWidth +
		param.BarWidth + param.PostBarWidth + 4 // + spaces
	buf.Grow(w)

	buf.WriteString(strings.Repeat(" ", param.PrePad))
	buf.WriteString(strutil.ResizeR(b.key, param.Post, param.KeyWidth))
	buf.WriteString(param.KeyDiv)
	buf.WriteRune(' ')
	buf.WriteString(strutil.ResizeR(b.msg, param.Post, param.MsgWidth))
	buf.WriteRune(' ')
	preBar := b.preBarF(b.curr, b.total, b.start, b.stopped)
	buf.WriteString(strutil.ResizeL(preBar, param.Post, param.PreBarWidth))
	buf.WriteRune(' ')
	b.makeBar(param, buf)
	buf.WriteRune(' ')
	postBar := b.postBarF(b.curr, b.total, b.start, b.stopped)
	buf.WriteString(strutil.ResizeL(postBar, param.Post, param.PostBarWidth))

	return buf.String()
}
