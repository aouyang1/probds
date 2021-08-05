package main

import (
	"github.com/OneOfOne/xxhash"
	"math/rand"
)

type Bloom struct {
	buckets   []byte
	hashFuncs []*xxhash.XXHash64
}

func NewBloom(bins, numHash int) *Bloom {
	b := new(Bloom)
	b.buckets = make([]byte, (bins-1)/8+1)
	b.hashFuncs = make([]*xxhash.XXHash64, 0, numHash)
	for i := 0; i < numHash; i++ {
		b.hashFuncs = append(b.hashFuncs, xxhash.NewS64(rand.Uint64()))
	}
	return b
}

func (b *Bloom) Reset() {
	for i := 0; i < len(b.buckets); i++ {
		b.buckets[i] = 0
	}
}

func (b *Bloom) Add(val string) {
	var bucketIdx, bucketOffset int

	for _, hf := range b.hashFuncs {
		bucketIdx, bucketOffset = b.hash(hf, []byte(val))
		b.buckets[bucketIdx] = b.buckets[bucketIdx] | bucketBitmask(bucketOffset)
	}
}

func (b *Bloom) Exists(val string) bool {
	var bucketIdx, bucketOffset int

	for _, hf := range b.hashFuncs {
		bucketIdx, bucketOffset = b.hash(hf, []byte(val))

		// got a bucket that was not set, so it is definitely not in the filter
		if b.buckets[bucketIdx]&bucketBitmask(bucketOffset) == 0 {
			return false
		}
	}
	return true
}

func (b *Bloom) hash(hf *xxhash.XXHash64, val []byte) (int, int) {
	hf.Write(val)
	idx := int(hf.Sum64() % uint64(len(b.buckets)*8))
	hf.Reset()
	return idx / 8, idx % 8
}

func bucketBitmask(offset int) byte {
	return (1 << (8 - offset - 1))
}
