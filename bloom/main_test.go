package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBloom(t *testing.T) {
	b := NewBloom(8*1024, 5)
	for i := 0; i < 1000; i++ {
		num := fmt.Sprintf("%d", rand.Uint64())
		b.Add(num)
		if !b.Exists(num) {
			t.Errorf("expected to find, %d, in bloom filter", i)
			break
		}
	}
	var tests int
	var hits int
	for i := 1; i < 10000; i++ {
		tests += 1
		num := fmt.Sprintf("%d", rand.Uint64())
		if b.Exists(num) {
			hits += 1
		}
	}
	t.Logf("false positive estimate: %.2f%%", float64(hits)/float64(tests)*100)
}

func BenchmarkBloom(b *testing.B) {
	numItems := 500000
	items := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fmt.Sprintf("%d", int(rand.Float64()*float64(numItems)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf := NewBloom(8*1024*1024, 5)
		for j := 0; j < len(items); j++ {
			bf.Add(items[j])
		}
	}
}

func BenchmarkBasicMap(b *testing.B) {
	numItems := 500000
	items := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fmt.Sprintf("%d", int(rand.Float64()*float64(numItems)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]struct{})
		for j := 0; j < len(items); j++ {
			m[items[j]] = struct{}{}
		}
	}
}
