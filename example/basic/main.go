package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/nu11ptr/cmpb"
)

const total = 100

var keys = []string{"bar1", "bar2"}

func main() {
	p := cmpb.New()

	for idx, key := range keys {
		b := p.NewBar(key+strconv.Itoa(idx), total)
		go func() {
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				b.Increment()
			}
		}()
	}

	p.Start()
	p.Wait()
}
