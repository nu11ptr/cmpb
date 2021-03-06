package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/nu11ptr/cmpb"
)

const total = 100

var (
	keys    = []string{"server1000", "server1001", "server1002"}
	actions = []string{"downloading...", "compiling source...", "fetching...", "committing work..."}
)

func main() {
	p := cmpb.New()

	for _, key := range keys {
		b := p.NewBar(key, total)
		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)

				if rand.Intn(total) == 1 {
					b.Stop("error!", "")
					break
				} else {
					action := actions[rand.Intn(len(actions))]
					b.SetMessage(action)
				}

				b.Increment()
			}
		}()
	}

	colors := cmpb.DefaultColors()
	colors.StopMsg = color.HiRedString
	p.SetColors(colors)
	p.Start()
	p.Wait()
}
