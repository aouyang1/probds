package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestComputeRegisterIndexAndLeftPos(t *testing.T) {
	testData := []struct {
		sum  uint32
		bits int
		reg  int
		pos  int
	}{
		{binary.BigEndian.Uint32([]byte{0x30, 0x01, 0x00, 0x00}), 4, 3, 12},
		{binary.BigEndian.Uint32([]byte{0xa0, 0x00, 0x00, 0xff}), 4, 10, 21},
	}

	for _, td := range testData {
		h := NewHyperLogLog(td.bits)
		reg, pos := h.computeRegisterIndexAndLeftPos(td.sum)
		if reg != td.reg {
			t.Errorf("expected register position %d, but got %d, for %+v", td.reg, reg, td)
		}
		if pos != td.pos {
			t.Errorf("expected left most one position %d, but got %d, for %+v", td.pos, pos, td)
		}
	}
}

func TestHyperLogLog(t *testing.T) {
	h := NewHyperLogLog(11)
	text := "as he crossed toward the pharmacy at the corner he involuntarily turned his head because of a burst of light that had ricocheted from his temple and saw with that quick smile with which we greet a rainbow or a rose a blindingly white parallelogram of sky being unloaded from the van a dresser with mirrors across which as across a cinema screen, passed a flawlessly clear reflection of boughs sliding and swaying not arboreally but with a human vacillation produced by the nature of those who were carrying this sky these boughs this gliding facade"
	parts := strings.Split(text, " ")
	basicMap := make(map[string]struct{})
	for _, p := range parts {
		h.Add(p)
		basicMap[p] = struct{}{}
	}
	t.Logf("estimating cardinality of %d, with actual of %d", h.Count(), len(basicMap))
}

func BenchmarkHyperLogLog(b *testing.B) {
	numItems := 10000000
	items := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fmt.Sprintf("%.9f", rand.Float64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := NewHyperLogLog(11)
		for j := 0; j < numItems; j++ {
			h.Add(items[j])
		}
		b.Logf("estimating %d unique items after inserting %d items", h.Count(), numItems)
	}
}

func BenchmarkBasicMap(b *testing.B) {
	numItems := 10000000
	items := make([]string, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fmt.Sprintf("%.9f", rand.Float64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		basicMap := make(map[string]struct{})
		for j := 0; j < numItems; j++ {
			basicMap[items[j]] = struct{}{}
		}
		b.Logf("estimating %d unique items after inserting %d items", len(basicMap), numItems)
	}
}
