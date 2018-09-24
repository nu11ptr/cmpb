package cmpb

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	defaultInterval = 200 * time.Millisecond
	defaultWidth    = 20 // Each char = 5%

	defaultLBracket = '['
	defaultRBracket = ']'
	defaultEmpty    = '-'
	defaultFull     = '='
	defaultCurr     = '>'

	slMapCap = 16
)

// Progress represents a collection of progress bars
type Progress struct {
	interval time.Duration
	width    int

	quitCh, waitCh chan struct{}
	wait           sync.WaitGroup
	out            io.Writer
	scrollUp       func(int, io.Writer)

	lBracket, rBracket, empty, full, curr rune

	bars   []*Bar
	barMap map[interface{}]*Bar
}

// New creates a new progress bar
func New() *Progress {
	return &Progress{
		interval: defaultInterval, width: defaultWidth,

		quitCh: make(chan struct{}), waitCh: make(chan struct{}),
		out: os.Stdout, scrollUp: ansiScrollUp,

		lBracket: defaultLBracket, rBracket: defaultRBracket, empty: defaultEmpty,
		full: defaultFull, curr: defaultCurr,

		bars: make([]*Bar, 0, slMapCap), barMap: make(map[interface{}]*Bar, slMapCap),
	}
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
		p.scrollUp(len(p.bars), p.out)
	}
	for _, bar := range p.bars {
		fmt.Fprintln(p.out, bar.String())
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
			case <-time.After(p.interval):
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
