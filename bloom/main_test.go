package main

import (
	"math/rand"
	"testing"
)

func TestBloom(t *testing.T) {
	b := NewBloom(8*1024, 5)
	for i := 0; i < 1000; i++ {
		num := rand.Uint64()
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
		num := rand.Uint64()
		if b.Exists(num) {
			hits += 1
		}
	}
	t.Logf("false positive estimate: %.2f%%", float64(hits)/float64(tests)*100)
}

func BenchmarkBloom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBloom(8*1024, 5)
	}
}
