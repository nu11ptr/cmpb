package cmpb

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nu11ptr/cmpb/strutil"
)

// BarColors represents a structure holding all bar drawing colors
type BarColors struct {
	Post, Key, KeyDiv, Msg, PreBar, LBracket, Empty, Full, Curr, RBracket,
	PostBar func(string, ...interface{}) string
}

// DefaultColors returns a set of default colors for rendering the bar
func DefaultColors() *BarColors {
	colors := new(BarColors)
	colors.SetAll(noOp)
	return colors
}

// SetAll sets all colors to the same color
func (b *BarColors) SetAll(f func(string, ...interface{}) string) {
	b.Post, b.Key, b.KeyDiv, b.Msg, b.PreBar, b.LBracket, b.Empty, b.Full, b.Curr, b.RBracket, b.PostBar =
		f, f, f, f, f, f, f, f, f, f, f
}

func noOp(s string, _ ...interface{}) string { return s }

// Bar represents a single progress bar
type Bar struct {
	key, msg            string
	curr, total         int
	lastRender, stopped bool
	start               time.Time
	preBarF, postBarF   func(int, int, time.Time, bool) string
	colors              BarColors

	p   *Progress
	mut sync.Mutex
}

func newBar(key string, total int, p *Progress) *Bar {
	return &Bar{
		key: key, msg: "", total: total, start: time.Now(), preBarF: CalcDur(), postBarF: CalcPct,
		colors: *DefaultColors(), p: p,
	}
}

// CalcDur calculates the duration since start time and returns a string
func CalcDur() func(int, int, time.Time, bool) string {
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
		return strutil.FmtDuration(dur)
	}
}

// CalcSteps calculates the steps completed so far and returns a string
func CalcSteps(curr, total int, start time.Time, stopped bool) string {
	return fmt.Sprintf("(%d/%d)", curr, total)
}

// CalcPct calculates the percentage of work complete and returns as a string
func CalcPct(curr, total int, start time.Time, stopped bool) string {
	pct := (curr * 100) / total
	return fmt.Sprintf("%d%%", pct)
}

// SetPreBar sets the prebar function decorator
func (b *Bar) SetPreBar(f func(int, int, time.Time, bool) string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.preBarF = f
}

// SetPostBar sets the postbar function decorator
func (b *Bar) SetPostBar(f func(int, int, time.Time, bool) string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.postBarF = f
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
			b.stopped = true
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

// Stop stops the updating of the bar and sets a final msg (if not ab empty string)
func (b *Bar) Stop(msg string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	if !b.stopped {
		b.stopped = true
		b.lastRender = true
		if msg != "" {
			b.msg = msg
		}
	}
}

// SetColors sets the colors used to render the bar
func (b *Bar) SetColors(colors *BarColors) {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.colors = *colors
}

// GetColors returns a copy of the colors used by this bar
func (b *Bar) GetColors() BarColors {
	b.mut.Lock()
	defer b.mut.Unlock()

	return b.colors
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

	// It can only be the last render a single time - subsequent calls return false
	lr := b.lastRender
	if lr {
		b.lastRender = false
	}
	return lr
}

func (b *Bar) makeBar(c *BarColors, param *Param, buf *bytes.Buffer) {
	buf.WriteString(c.LBracket(string(param.LBracket)))

	full := b.curr * (param.BarWidth - 2) / b.total
	empty := (param.BarWidth - 2) - full
	if full > 0 {
		if empty > 0 {
			full--
		}
		buf.WriteString(c.Full(strings.Repeat(string(param.Full), full)))
		if empty > 0 {
			buf.WriteString(c.Curr(string(param.Curr)))
		}
	}
	buf.WriteString(c.Empty(strings.Repeat(string(param.Empty), empty)))

	buf.WriteString(c.RBracket(string(param.RBracket)))
}

func (b *Bar) String() string {
	b.mut.Lock()
	defer b.mut.Unlock()

	buf := new(bytes.Buffer)
	param := &b.p.param
	w := param.PrePad + param.KeyWidth + param.MsgWidth + param.PreBarWidth +
		param.BarWidth + param.PostBarWidth + 5 // spaces + keyDiv
	buf.Grow(w)
	c := &b.colors

	buf.WriteString(strings.Repeat(" ", param.PrePad))
	buf.WriteString(strutil.ResizeR(c.Key(b.key), c.Post(param.Post), param.KeyWidth))
	buf.WriteString(c.KeyDiv(string(param.KeyDiv)))
	buf.WriteRune(' ')
	buf.WriteString(strutil.ResizeR(c.Msg(b.msg), c.Post(param.Post), param.MsgWidth))
	buf.WriteRune(' ')
	preBar := c.PreBar(b.preBarF(b.curr, b.total, b.start, b.stopped))
	buf.WriteString(strutil.ResizeL(preBar, c.Post(param.Post), param.PreBarWidth))
	buf.WriteRune(' ')
	b.makeBar(c, param, buf)
	buf.WriteRune(' ')
	postBar := c.PostBar(b.postBarF(b.curr, b.total, b.start, b.stopped))
	buf.WriteString(strutil.ResizeL(postBar, c.Post(param.Post), param.PostBarWidth))

	return buf.String()
}
