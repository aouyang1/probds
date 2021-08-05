package main

import (
	"github.com/OneOfOne/xxhash"
	"math"
	"math/rand"
)

type CountMinSketch struct {
	numHashFuncs int
	numBuckets   int
	hashBuckets  [][]int
	hashFuncs    []*xxhash.XXHash64
}

func NewCountMinSketch(bins, numHash int) *CountMinSketch {
	c := new(CountMinSketch)
	c.numHashFuncs = numHash
	c.numBuckets = bins
	c.hashBuckets = make([][]int, numHash)
	for i := 0; i < numHash; i++ {
		c.hashBuckets[i] = make([]int, bins)
	}

	c.hashFuncs = make([]*xxhash.XXHash64, 0, numHash)
	for i := 0; i < numHash; i++ {
		c.hashFuncs = append(c.hashFuncs, xxhash.NewS64(rand.Uint64()))
	}

	return c
}

func (c *CountMinSketch) Reset() {
	for i := 0; i < len(c.hashBuckets); i++ {
		for j := 0; j < len(c.hashBuckets[j]); j++ {
			c.hashBuckets[i][j] = 0
		}
	}
}

func (c *CountMinSketch) Add(val string) {
	var bucketIdx int

	for i, hf := range c.hashFuncs {
		bucketIdx = c.hash(hf, []byte(val))
		c.hashBuckets[i][bucketIdx] += 1
	}
}

func (c *CountMinSketch) Count(val string) int {
	var bucketIdx int
	var currVal int
	minVal := math.MaxInt64

	for i, hf := range c.hashFuncs {
		bucketIdx = c.hash(hf, []byte(val))
		currVal = c.hashBuckets[i][bucketIdx]
		if currVal < minVal {
			minVal = c.hashBuckets[i][bucketIdx]

			// means that the current value hasn't ever been set, so stop searching
			// with remaining hash functions
			if minVal == 0 {
				break
			}
		}
	}
	return minVal
}

func (c *CountMinSketch) hash(hf *xxhash.XXHash64, val []byte) int {
	hf.Write(val)
	idx := int(hf.Sum64() % uint64(c.numBuckets))
	hf.Reset()
	return idx
}
