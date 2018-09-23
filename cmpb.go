package cmpb

import "time"

const (
	defaultInterval = 200 * time.Millisecond
	defaultWidth    = 20 // Each char = 5%
)

// Progress represents a collection of progress bars
type Progress struct {
	interval time.Duration
	width    int
	quitCh   chan struct{}

	bars []*Bar
}

// Bar represents a single progress bar
type Bar struct {
}

// New creates a new progress bar
func New() *Progress {
	return &Progress{interval: defaultInterval, width: defaultWidth, quitCh: make(chan struct{})}
}

// NewBar creates a new progress bar and adds it to the progress bar collection
func (p *Progress) NewBar() *Bar {
	b := new(Bar)
	p.bars = append(p.bars, b)
	return b
}

func (p *Progress) render() {
}

// Start begins rendering of the progress bars
func (p *Progress) Start() {
	go func() {
		for {
			select {
			case <-time.After(p.interval):
				p.render()
			case <-p.quitCh:
				break
			}
		}
	}()
}

// Stop stops the render of the progress bars
func (p *Progress) Stop() {
	close(p.quitCh)
}
