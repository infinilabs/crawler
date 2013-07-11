package bloom

import (
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"github.com/pmylund/go-bitset"
)

type filter struct {
	m uint32
	k uint32
	h hash.Hash64
}

func (f *filter) bits(data []byte) []uint32 {
	f.h.Reset()
	f.h.Write(data)
	d := f.h.Sum(nil)
	a := binary.BigEndian.Uint32(d[4:8])
	b := binary.BigEndian.Uint32(d[0:4])
	is := make([]uint32, f.k)
	for i := uint32(0); i < f.k; i++ {
		is[i] = (a + b*i) % f.m
	}
	return is
}

func newFilter(m, k uint32) *filter {
	return &filter{
		m: m,
		k: k,
		h: fnv.New64(),
	}
}

func estimates(n uint32, p float64) (uint32, uint32) {
	nf := float64(n)
	log2 := math.Log(2)
	m := -1 * nf * math.Log(p) / math.Pow(log2, 2)
	k := math.Ceil(log2 * m / nf)

	words := m + 31>>5
	if words >= math.MaxInt32 {
		panic(fmt.Sprintf("A 32-bit bloom filter with n %d and p %f requires a 32-bit bitset with a slice of %f words, but slices cannot contain more than %d elements. Please use the equivalent 64-bit bloom filter, e.g. New64(), instead.", n, p, words, math.MaxInt32-1))
	} else if m > math.MaxUint32 {
		panic(fmt.Sprintf("A 32-bit bloom filter with n %d and p %f requires a 32-bit bitset with %d bits, but this number overflows an uint32. Please use the equivalent 64-bit bloom filter, e.g. New64(), instead.", n, p, m))
	}
	return uint32(m), uint32(k)
}

// A standard bloom filter using the 64-bit FNV-1a hash function.
type Filter struct {
	*filter
	b *bitset.Bitset32
}

// Check whether data was previously added to the filter. Returns true if
// yes, with a false positive chance near the ratio specified upon creation
// of the filter. The result cannot be falsely negative.
func (f *Filter) Test(data []byte) bool {
	for _, i := range f.bits(data) {
		if !f.b.Test(i) {
			return false
		}
	}
	return true
}

// Add data to the filter.
func (f *Filter) Add(data []byte) {
	for _, i := range f.bits(data) {
		f.b.Set(i)
	}
}

// Resets the filter.
func (f *Filter) Reset() {
	f.b.Reset()
}

// Create a bloom filter with an expected n number of items, and an acceptable
// false positive rate of p, e.g. 0.01.
func New(n int, p float64) *Filter {
	m, k := estimates(uint32(n), p)
	f := &Filter{
		newFilter(m, k),
		bitset.New32(m),
	}
	return f
}

// A counting bloom filter using the 64-bit FNV-1a hash function. Supports
// removing items from the filter.
type CountingFilter struct {
	*filter
	b []*bitset.Bitset32
}

// Checks whether data was previously added to the filter. Returns true if
// yes, with a false positive chance near the ratio specified upon creation
// of the filter. The result cannot cannot be falsely negative (unless one
// has removed an item that wasn't actually added to the filter previously.)
func (f *CountingFilter) Test(data []byte) bool {
	b := f.b[0]
	for _, v := range f.bits(data) {
		if !b.Test(v) {
			return false
		}
	}
	return true
}

// Adds data to the filter.
func (f *CountingFilter) Add(data []byte) {
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
			nb := bitset.New32(f.b[0].Len())
			f.b = append(f.b, nb)
			nb.Set(v)
		}
	}
}

// Removes data from the filter. This exact data must have been previously added
// to the filter, or future results will be inconsistent.
func (f *CountingFilter) Remove(data []byte) {
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
func (f *CountingFilter) Reset() {
	f.b = f.b[:1]
	f.b[0].Reset()
}

// Create a counting bloom filter with an expected n number of items, and an
// acceptable false positive rate of p. Counting bloom filters support
// the removal of items from the filter.
func NewCounting(n int, p float64) *CountingFilter {
	m, k := estimates(uint32(n), p)
	f := &CountingFilter{
		newFilter(m, k),
		[]*bitset.Bitset32{bitset.New32(m)},
	}
	return f
}

// A layered bloom filter using the 64-bit FNV-1a hash function.
type LayeredFilter struct {
	*filter
	b []*bitset.Bitset32
}

// Checks whether data was previously added to the filter. Returns the number of
// the last layer where the data was added, e.g. 1 for the first layer, and a
// boolean indicating whether the data was added to the filter at all. The check
// has a false positive chance near the ratio specified upon creation of the
// filter. The result cannot be falsely negative.
func (f *LayeredFilter) Test(data []byte) (int, bool) {
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
func (f *LayeredFilter) Add(data []byte) int {
	is := f.bits(data)
	var (
		i int
		v *bitset.Bitset32
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
	nb := bitset.New32(f.b[0].Len())
	f.b = append(f.b, nb)
	for _, v := range is {
		nb.Set(v)
	}
	return i + 2
}

// Resets the filter.
func (f *LayeredFilter) Reset() {
	f.b = f.b[:1]
	f.b[0].Reset()
}

// Create a layered bloom filter with an expected n number of items, and an
// acceptable false positive rate of p. Layered bloom filters can be used
// to keep track of a certain, arbitrary count of items, e.g. to check if some
// given data was added to the filter 10 times or less.
func NewLayered(n int, p float64) *LayeredFilter {
	m, k := estimates(uint32(n), p)
	f := &LayeredFilter{
		newFilter(m, k),
		[]*bitset.Bitset32{bitset.New32(m)},
	}
	return f
}
