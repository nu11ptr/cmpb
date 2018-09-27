package cmpb_test

import (
	"testing"
	"time"

	"github.com/nu11ptr/cmpb"
)

func TestCalcPct(t *testing.T) {
	tests := []struct {
		name        string
		curr, total int
		output      string
	}{
		{"0", 0, 10, "0%"},
		{"Round", 33, 99, "33%"},
		{"100", 10, 10, "100%"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := cmpb.CalcPct(test.curr, test.total, time.Now(), false)
			if output != test.output {
				t.Error("want", test.output, "got", output)
			}
		})
	}
}

func TestSteps(t *testing.T) {
	output := cmpb.CalcSteps(10, 33, time.Now(), false)
	if output != "(10/33)" {
		t.Error("want", "(10/33)", "got", output)
	}
}

func TestString(t *testing.T) {
	p := cmpb.New()
	b := p.NewBar("bar", 10)

	t.Run("First", func(t *testing.T) {
		expected := "bar       :                               0s [--------------------]   0%"
		output := b.String()
		if output != expected {
			t.Error("want", expected, "got", output)
		}
	})
	t.Run("Increment", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			b.Increment()
		}
		expected := "bar       :                               0s [=========>----------]  50%"
		output := b.String()
		if output != expected {
			t.Error("want", expected, "got", output)
		}
	})
	t.Run("Complete", func(t *testing.T) {
		b.Update(10)
		b.SetMessage("done")
		expected := "bar       : done                          0s [====================] 100%"
		output := b.String()
		if output != expected {
			t.Error("want", expected, "got", output)
		}
	})
}
