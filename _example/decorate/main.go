package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nu11ptr/cmpb"
)

const total = 100

var (
	keys    = []string{"server1000", "server1001", "server1002"}
	actions = []string{"downloading...", "compiling source...", "fetching...", "committing work..."}
)

func calcStepsDur() func(int, int, time.Time, bool) string {
	f := cmpb.CalcDur()

	return func(curr, total int, start time.Time, stopped bool) string {
		return fmt.Sprintf("%s %s", cmpb.CalcSteps(curr, total, start, stopped),
			f(curr, total, start, stopped))
	}
}

func main() {
	param := cmpb.DefaultParam()
	param.PreBarWidth = 13
	p := cmpb.NewWithParam(param)

	for _, key := range keys {
		b := p.NewBar(key, total)
		b.SetPreBar(calcStepsDur())

		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
				action := actions[rand.Intn(len(actions))]
				b.SetMessage(action)
				b.Increment()
			}
		}()
	}

	p.Start()
	p.Wait()
}
