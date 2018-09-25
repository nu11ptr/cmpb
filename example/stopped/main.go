package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/nu11ptr/cmpb"
)

const total = 100

var (
	keys    = []string{"server100", "server101", "server102"}
	actions = []string{"downloading...", "compiling source...", "fetching...", "committing work..."}
)

func main() {
	p := cmpb.New()

	for _, key := range keys {
		b := p.NewBar(color.HiBlueString(key), total)
		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
				action := actions[rand.Intn(len(actions))]

				if rand.Intn(total) == 1 {
					b.Stop(color.HiRedString("error!"))
					break
				} else {
					b.SetMessage(color.HiYellowString(action))
				}

				b.Increment()
			}
		}()
	}

	p.Start()
	p.Wait()
}
