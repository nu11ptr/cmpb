# cmpb
[![Build Status](https://travis-ci.org/nu11ptr/cmpb.svg?branch=master)](https://travis-ci.org/nu11ptr/cmpb)
[![Build status](https://ci.appveyor.com/api/projects/status/2kxwqb49ihfvaiy3/branch/master?svg=true)](https://ci.appveyor.com/project/nu11ptr/cmpb/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/nu11ptr/cmpb/badge.svg?branch=master)](https://coveralls.io/github/nu11ptr/cmpb?branch=master)
[![Maintainability](https://api.codeclimate.com/v1/badges/253b30e054c6844f3e9c/maintainability)](https://codeclimate.com/github/nu11ptr/cmpb/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/nu11ptr/cmpb)](https://goreportcard.com/report/github.com/nu11ptr/cmpb)
[![GoDoc](https://godoc.org/github.com/nu11ptr/cmpb?status.svg)](https://godoc.org/github.com/nu11ptr/cmpb)

A color multi-progress bar for Go terminal applications. Works on ANSI compatible terminals... and also on Windows!

# Demo
![Image](demo.svg)

# Demo Source

```go
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
	colors.Key = color.HiBlueString
	colors.Msg, colors.Empty = color.HiYellowString, color.HiYellowString
	colors.Full = color.HiGreenString
	colors.Curr = color.GreenString
	colors.PreBar, colors.PostBar = color.HiMagentaString, color.HiMagentaString

	for _, key := range keys {
		b := p.NewBar(key, total)
		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
				action := actions[rand.Intn(len(actions))]
				b.SetMessage(action)
				b.Increment()
			}
		}()
	}

	p.SetColors(colors)
	p.Start()
	p.Wait()
}
```