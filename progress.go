package cmpb

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	defaultInterval = 200 * time.Millisecond
	defaultWidth    = 20 // Each char = 5%

	slMapCap = 16
)

var (
	defaultLBracket = color.HiCyanString("[")
	defaultRBracket = color.HiCyanString("]")
	defaultEmpty    = color.HiYellowString("-")
	defaultFull     = color.HiGreenString("=")
	defaultCurr     = color.HiGreenString(">")
)

// Param represents the parameters for a Progress
type Param struct {
	Interval time.Duration
	Width    int
	Out      io.Writer
	ScrollUp func(int, io.Writer)

	LBracket, RBracket, Empty, Full, Curr string
}

// DefaultParam builds a Param struct with default values
func DefaultParam() *Param {
	return &Param{
		Interval: defaultInterval, Width: defaultWidth,
		Out: color.Output, ScrollUp: ansiScrollUp,

		LBracket: defaultLBracket, RBracket: defaultRBracket, Empty: defaultEmpty,
		Full: defaultFull, Curr: defaultCurr,
	}
}

// Progress represents a collection of progress bars
type Progress struct {
	param Param

	quitCh, waitCh chan struct{}
	wait           sync.WaitGroup

	bars   []*Bar
	barMap map[interface{}]*Bar
}

// NewWithParam creates a new progress bar collection with specified params
func NewWithParam(param *Param) *Progress {
	return &Progress{
		param:  *param,
		quitCh: make(chan struct{}), waitCh: make(chan struct{}),
		bars: make([]*Bar, 0, slMapCap), barMap: make(map[interface{}]*Bar, slMapCap),
	}
}

// New creates a new progress bar collection with default params
func New() *Progress {
	return NewWithParam(DefaultParam())
}

func ansiScrollUp(rows int, out io.Writer) {
	fmt.Fprintf(out, "\x1b[%dA", rows)
}

// NewBar creates a new progress bar and adds it to the progress bar collection
func (p *Progress) NewBar(key interface{}, total int) *Bar {
	b := newBar(p, total)
	p.bars = append(p.bars, b)
	p.barMap[key] = b
	p.wait.Add(1)
	return b
}

// Bar returns the bar stored the given key. The value is nil if it can't be found
func (p *Progress) Bar(key interface{}) *Bar {
	b, _ := p.barMap[key]
	return b
}

func (p *Progress) render(scrollUp bool) {
	if scrollUp {
		p.param.ScrollUp(len(p.bars), p.param.Out)
	}
	for _, bar := range p.bars {
		fmt.Fprintln(p.param.Out, bar.String())
	}
	// Done as 2nd pass so all bars are always rendered
	for _, bar := range p.bars {
		if bar.lastRender {
			bar.lastRender = false
			p.wait.Done()
		}
	}
}

// Start begins rendering of the progress bars
func (p *Progress) Start() {
	go func() {
		firstTime := true
		for {
			select {
			case <-time.After(p.param.Interval):
				p.render(!firstTime)
				firstTime = false
			case <-p.quitCh:
				break
			}
		}
	}()
}

// Stop stops the render of the progress bars
func (p *Progress) Stop() {
	select {
	case <-p.quitCh:
	default:
		close(p.quitCh)
	}
}

// Wait waits for progress to be finished or cancelled. It can only be called once
func (p *Progress) Wait() {
	go func() {
		p.wait.Wait()
		close(p.waitCh)
	}()

	select {
	case <-p.quitCh:
	case <-p.waitCh:
	}
}
