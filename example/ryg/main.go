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
	param := cmpb.DefaultParam()
	param.InlineExtMsg = false
	p := cmpb.NewWithParam(param)

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
					b.Stop("error!", "A massive error occurred:\n    The error is catastrophic and cannot be recovered from")
					break
				} else {
					if i == total-1 {
						colors.SetAll(color.HiGreenString)
						b.SetColors(colors)
						b.SetMessage("completed")
					} else {
						action := actions[rand.Intn(len(actions))]
						b.SetMessage(action)
					}
				}

				b.Increment()
			}
		}()
	}

	p.SetColors(colors)
	p.Start()
	p.Wait()
}
