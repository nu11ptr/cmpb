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

	colors := new(cmpb.BarColors)
	colors.SetAll(color.HiYellowString)

	for _, key := range keys {
		b := p.NewBar(key, total)
		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)

				if rand.Intn(total*3) == 1 {
					colors.SetAll(color.HiRedString)
					b.SetColors(colors)
					b.Stop("error!", "A massive error occurred:\n\t\tThe error is catastrophic and cannot be recovered from")
					break
				} else {
					action := actions[rand.Intn(len(actions))]
					if i == total-1 {
						colors.SetAll(color.HiGreenString)
						b.SetColors(colors)
					}
					b.SetMessage(action)
				}

				b.Increment()
			}
		}()
	}

	p.SetColors(colors)
	p.Start()
	p.Wait()
}
