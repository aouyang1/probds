package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestCountMinSketch(t *testing.T) {
	bins := 1024
	hashFuncs := 3
	c := NewCountMinSketch(bins, hashFuncs)
	text := "as he crossed toward the pharmacy at the corner he involuntarily turned his head because of a burst of light that had ricocheted from his temple and saw with that quick smile with which we greet a rainbow or a rose a blindingly white parallelogram of sky being unloaded from the van a dresser with mirrors across which as across a cinema screen, passed a flawlessly clear reflection of boughs sliding and swaying not arboreally but with a human vacillation produced by the nature of those who were carrying this sky these boughs this gliding facade"
	parts := strings.Split(text, " ")
	for _, p := range parts {
		c.Add(p)
	}

	cases := []struct {
		word  string
		count int
	}{
		{"a", 8},
		{"of", 5},
		{"with", 4},
		{"the", 4},
		{"which", 2},
	}
	for _, d := range cases {
		res := c.Count(d.word)
		if res != d.count {
			t.Errorf("expected %d for %q, but got %d", d.count, d.word, res)
		}
	}
}

func BenchmarkCountMinSketch(b *testing.B) {
	numItems := 10000000
	maxNum := 1000000.0
	items := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fmt.Sprintf("%d", int(rand.Float64()*maxNum))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bins := 1024 * 1024
		hashFuncs := 3
		c := NewCountMinSketch(bins, hashFuncs)
		for j := 0; j < numItems; j++ {
			c.Add(items[j])
		}
		b.Logf("estimating %d count for %s", c.Count(fmt.Sprintf("1")), fmt.Sprintf("1"))
	}
}

func BenchmarkBasicMap(b *testing.B) {
	numItems := 10000000
	maxNum := 1000000.0
	items := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fmt.Sprintf("%d", int(rand.Float64()*maxNum))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		basicMap := make(map[string]int)
		for j := 0; j < numItems; j++ {
			basicMap[items[j]] += 1
		}
		b.Logf("estimating %d count for %s", basicMap[fmt.Sprintf("1")], fmt.Sprintf("1"))
	}
}
