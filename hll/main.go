package main

import (
	"github.com/OneOfOne/xxhash"
	"math"
)

type HyperLogLog struct {
	bits      int
	registers []int
	hashFunc  *xxhash.XXHash32
	alpha     float64
}

func NewHyperLogLog(bits int) *HyperLogLog {
	h := new(HyperLogLog)
	h.bits = bits
	numRegisters := int(math.Pow(2.0, float64(bits)))
	h.registers = make([]int, numRegisters)
	h.hashFunc = xxhash.New32()
	h.alpha = h.computeAlpha()
	return h
}

func (h *HyperLogLog) Add(val string) {
	reg, pos := h.hash([]byte(val))
	if pos > h.registers[reg] {
		h.registers[reg] = pos
	}
}

func (h *HyperLogLog) Count() int {
	var z float64
	var v int
	for _, val := range h.registers {
		if val == 0 {
			v += 1
		}
		z += math.Pow(2.0, -float64(val))
	}
	z = 1.0 / z
	m := float64(len(h.registers))
	estimate := h.alpha * m * m * z
	if estimate <= 5.0/2.0*m && v != 0 {
		return int(m * math.Log(m/float64(v)))
	}
	if estimate > 1.0/30.0*math.Pow(2.0, 32.0) {
		return int(-math.Pow(2.0, 32.0) * math.Log(1.0-estimate/math.Pow(2.0, 32.0)))
	}
	return int(estimate)
}

//computeAlpha returns the bias factor based on the number of registers expected to be
// a power of 2. Lowest appropriate register length is 16.
func (h *HyperLogLog) computeAlpha() float64 {
	m := float64(len(h.registers))
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/m)
	}
}

// hash compute the hash of the input value and returns the register it belongs to along with
// the position of the left most 1 in the hash after the register bits.
func (h *HyperLogLog) hash(val []byte) (int, int) {
	h.hashFunc.Write(val)
	sum := h.hashFunc.Sum32()
	h.hashFunc.Reset()
	return h.computeRegisterIndexAndLeftPos(sum)
}

// computeRegisterIndexAndLeftPos returns the register index based on a 64 bit hash and
// the HLL configured bits along with the position of the left most bit set to 1 after the
// register based bits.
// e.g. bits = 4
// 0101|0010111100...
// register index:      5
// left most 1 postion: 3
func (h *HyperLogLog) computeRegisterIndexAndLeftPos(sum uint32) (int, int) {
	var pos int
	regIdx := int(sum >> (32 - h.bits))
	masked := sum &^ ((uint32(math.Pow(2.0, float64(h.bits))) - 1) << (32 - h.bits))
	for i := 0; i < 32-h.bits; i++ {
		if masked == 1 {
			pos = 32 - h.bits - i
			break
		}
		masked = masked >> 1
	}
	return regIdx, pos
}
