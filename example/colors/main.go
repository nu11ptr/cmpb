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
	colors.Post, colors.KeyDiv, colors.LBracket, colors.RBracket =
		color.HiCyanString, color.HiCyanString, color.HiCyanString, color.HiCyanString
	colors.Empty = color.HiYellowString
	colors.Full = color.HiGreenString
	colors.Curr = color.GreenString
	colors.PreBar, colors.PostBar = color.HiMagentaString, color.HiMagentaString

	for _, key := range keys {
		b := p.NewBar(color.HiBlueString(key), total)
		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
				action := actions[rand.Intn(len(actions))]
				b.SetMessage(color.HiYellowString(action))
				b.Increment()
			}
		}()
	}

	p.SetColors(colors)
	p.Start()
	p.Wait()
}
