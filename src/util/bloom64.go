package bloom

import (
	"hash"
	"hash/crc64"
	"hash/fnv"
	"math"
	"util/bitset"
)

type filter64 struct {
	m  uint64
	k  uint64
	h  hash.Hash64
	oh hash.Hash64
}

func (f *filter64) bits(data []byte) []uint64 {
	f.h.Reset()
	f.h.Write(data)
	a := f.h.Sum64()

	f.oh.Reset()
	f.oh.Write(data)
	b := f.oh.Sum64()

	is := make([]uint64, f.k)
	for i := uint64(0); i < f.k; i++ {
		is[i] = (a + b*i) % f.m
	}
	return is
}

func newFilter64(m, k uint64) *filter64 {
	return &filter64{
		m:  m,
		k:  k,
		h:  fnv.New64(),
		oh: crc64.New(crc64.MakeTable(crc64.ECMA)),
	}
}

func estimates64(n uint64, p float64) (uint64, uint64) {
	nf := float64(n)
	log2 := math.Log(2)
	m := -1 * nf * math.Log(p) / math.Pow(log2, 2)
	k := math.Ceil(log2 * m / nf)
	return uint64(m), uint64(k)
}

// A standard 64-bit bloom filter using the 64-bit FNV-1a hash function.
type Filter64 struct {
	*filter64
	b *bitset.Bitset64
}

// Check whether data was previously added to the filter. Returns true if
// yes, with a false positive chance near the ratio specified upon creation
// of the filter. The result cannot be falsely negative.
func (f *Filter64) Test(data []byte) bool {
	for _, i := range f.bits(data) {
		if !f.b.Test(i) {
			return false
		}
	}
	return true
}

// Add data to the filter.
func (f *Filter64) Add(data []byte) {
	for _, i := range f.bits(data) {
		f.b.Set(i)
	}
}

// Resets the filter.
func (f *Filter64) Reset() {
	f.b.Reset()
}

// Create a bloom filter with an expected n number of items, and an acceptable
// false positive rate of p, e.g. 0.01 for 1%.
func New64(n int64, p float64) *Filter64 {
	m, k := estimates64(uint64(n), p)
	f := &Filter64{
		newFilter64(m, k),
		bitset.New64(m),
	}
	return f
}

// A counting bloom filter using the 64-bit FNV-1a hash function. Supports
// removing items from the filter.
type CountingFilter64 struct {
	*filter64
	b []*bitset.Bitset64
}

// Checks whether data was previously added to the filter. Returns true if
// yes, with a false positive chance near the ratio specified upon creation
// of the filter. The result cannot cannot be falsely negative (unless one
// has removed an item that wasn't actually added to the filter previously.)
func (f *CountingFilter64) Test(data []byte) bool {
	b := f.b[0]
	for _, v := range f.bits(data) {
		if !b.Test(v) {
			return false
		}
	}
	return true
}

// Adds data to the filter.
func (f *CountingFilter64) Add(data []byte) {
	for _, v := range f.bits(data) {
		done := false
		for _, ov := range f.b {
			if !ov.Test(v) {
				done = true
				ov.Set(v)
				break
			}
		}
		if !done {
			nb := bitset.New64(f.b[0].Len())
			f.b = append(f.b, nb)
			nb.Set(v)
		}
	}
}

// Removes data from the filter. This exact data must have been previously added
// to the filter, or future results will be inconsistent.
func (f *CountingFilter64) Remove(data []byte) {
	last := len(f.b) - 1
	for _, v := range f.bits(data) {
		for oi := last; oi >= 0; oi-- {
			ov := f.b[oi]
			if ov.Test(v) {
				ov.Clear(v)
				break
			}
		}
	}
}

// Resets the filter.
func (f *CountingFilter64) Reset() {
	f.b = f.b[:1]
	f.b[0].Reset()
}

// Create a counting bloom filter with an expected n number of items, and an
// acceptable false positive rate of p. Counting bloom filters support
// the removal of items from the filter.
func NewCounting64(n int64, p float64) *CountingFilter64 {
	m, k := estimates64(uint64(n), p)
	f := &CountingFilter64{
		newFilter64(m, k),
		[]*bitset.Bitset64{bitset.New64(m)},
	}
	return f
}

// A layered bloom filter using the 64-bit FNV-1a hash function.
type LayeredFilter64 struct {
	*filter64
	b []*bitset.Bitset64
}

// Checks whether data was previously added to the filter. Returns the number of
// the last layer where the data was added, e.g. 1 for the first layer, and a
// boolean indicating whether the data was added to the filter at all. The check
// has a false positive chance near the ratio specified upon creation of the
// filter. The result cannot be falsely negative.
func (f *LayeredFilter64) Test(data []byte) (int, bool) {
	is := f.bits(data)
	for i := len(f.b) - 1; i >= 0; i-- {
		v := f.b[i]
		last := len(is) - 1
		for oi, ov := range is {
			if !v.Test(ov) {
				break
			}
			if oi == last {
				// Every test was positive at this layer
				return i + 1, true
			}
		}
	}
	return 0, false
}

// Adds data to the filter. Returns the number of the layer where the data
// was added, e.g. 1 for the first layer.
func (f *LayeredFilter64) Add(data []byte) int {
	is := f.bits(data)
	var (
		i int
		v *bitset.Bitset64
	)
	for i, v = range f.b {
		here := false
		for _, ov := range is {
			if here {
				v.Set(ov)
			} else if !v.Test(ov) {
				here = true
				v.Set(ov)
			}
		}
		if here {
			return i + 1
		}
	}
	nb := bitset.New64(f.b[0].Len())
	f.b = append(f.b, nb)
	for _, v := range is {
		nb.Set(v)
	}
	return i + 2
}

// Resets the filter.
func (f *LayeredFilter64) Reset() {
	f.b = f.b[:1]
	f.b[0].Reset()
}

// Create a layered bloom filter with an expected n number of items, and an
// acceptable false positive rate of p. Layered bloom filters can be used
// to keep track of a certain, arbitrary count of items, e.g. to check if some
// given data was added to the filter 10 times or less.
func NewLayered64(n int64, p float64) *LayeredFilter64 {
	m, k := estimates64(uint64(n), p)
	f := &LayeredFilter64{
		newFilter64(m, k),
		[]*bitset.Bitset64{bitset.New64(m)},
	}
	return f
}
