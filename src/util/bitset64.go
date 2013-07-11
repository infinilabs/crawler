package bitset

import (
	"bytes"
	"fmt"
	"math"
)

const (
	sw_64   uint64 = 64
	slg2_64 uint64 = 6
	m1_64   uint64 = 0x5555555555555555 // 0101...
	m2_64   uint64 = 0x3333333333333333 // 00110011..
	m4_64   uint64 = 0x0f0f0f0f0f0f0f0f // 00001111...
	hff_64  uint64 = 0xffffffffffffffff // all ones
)

func wordsNeeded64(n uint64) uint64 {
	if n == 0 {
		return 1
	} else if n > math.MaxUint64-sw_64 {
		return math.MaxUint64 >> slg2_64
	}
	return (n + (sw_64 - 1)) >> slg2_64
}

type Bitset64 struct {
	n uint64
	b []uint64
}

// Returns the current size of the bitset.
func (b *Bitset64) Len() uint64 {
	return b.n
}

// Test whether bit i is set.
func (b *Bitset64) Test(i uint64) bool {
	if i >= b.n {
		return false
	}
	return ((b.b[i>>slg2_64] & (1 << (i & (sw_64 - 1)))) != 0)
}

// Set bit i to 1.
func (b *Bitset64) Set(i uint64) {
	if i >= b.n {
		nsize := wordsNeeded64(i + 1)
		l := uint64(len(b.b))
		if nsize > l {
			nb := make([]uint64, nsize-l)
			b.b = append(b.b, nb...)
		}
		b.n = i + 1
	}
	b.b[i>>slg2_64] |= (1 << (i & (sw_64 - 1)))
}

// Set bit i to 0.
func (b *Bitset64) Clear(i uint64) {
	if i >= b.n {
		return
	}
	b.b[i>>slg2_64] &^= 1 << (i & (sw_64 - 1))
}

// Flip bit i.
func (b *Bitset64) Flip(i uint64) {
	if i >= b.n {
		b.Set(i)
	}
	b.b[i>>slg2_64] ^= 1 << (i & (sw_64 - 1))
}

// Clear all bits in the bitset.
func (b *Bitset64) Reset() {
	for i := range b.b {
		b.b[i] = 0
	}
}

// Get the number of words used in the bitset.
func (b *Bitset64) wordCount() uint64 {
	return wordsNeeded64(b.n)
}

// Clone the bitset.
func (b *Bitset64) Clone() *Bitset64 {
	c := New64(b.n)
	copy(c.b, b.b)
	return c
}

// Copy the bitset into another bitset, returning the size of the destination
// bitset.
func (b *Bitset64) Copy(c *Bitset64) (n uint64) {
	copy(c.b, b.b)
	n = c.n
	if b.n < c.n {
		n = b.n
	}
	return
}

func popCountUint64(x uint64) uint64 {
	x -= (x >> 1) & m1_64                // put count of each 2 bits into those 2 bits
	x = (x & m2_64) + ((x >> 2) & m2_64) // put count of each 4 bits into those 4 bits 
	x = (x + (x >> 4)) & m4_64           // put count of each 8 bits into those 8 bits 
	x += x >> 8                          // put count of each 16 bits into their lowest 8 bits
	x += x >> 16                         // put count of each 32 bits into their lowest 8 bits
	x += x >> 32                         // put count of each 64 bits into their lowest 8 bits
	return x & 0x7f
}

// Get the number of set bits in the bitset.
func (b *Bitset64) Count() uint64 {
	sum := uint64(0)
	for _, w := range b.b {
		sum += popCountUint64(w)
	}
	return sum
}

// Test if two bitsets are equal. Returns true if both bitsets are the same
// size and all the same bits are set in both bitsets.
func (b *Bitset64) Equal(c *Bitset64) bool {
	if b.n != c.n {
		return false
	}
	for p, v := range b.b {
		if c.b[p] != v {
			return false
		}
	}
	return true
}

// Bitset &^ (and or); difference between receiver and another set.
func (b *Bitset64) Difference(ob *Bitset64) (result *Bitset64) {
	result = b.Clone() // clone b (in case b is bigger than ob)
	szl := ob.wordCount()
	l := uint64(len(b.b))
	for i := uint64(0); i < l; i++ {
		if i >= szl {
			break
		}
		result.b[i] = b.b[i] &^ ob.b[i]
	}
	return
}

func sortByLength64(a *Bitset64, b *Bitset64) (ap *Bitset64, bp *Bitset64) {
	if a.n <= b.n {
		ap, bp = a, b
	} else {
		ap, bp = b, a
	}
	return
}

// Bitset & (and); intersection of receiver and another set.
func (b *Bitset64) Intersection(ob *Bitset64) (result *Bitset64) {
	b, ob = sortByLength64(b, ob)
	result = New64(b.n)
	for i, w := range b.b {
		result.b[i] = w & ob.b[i]
	}
	return
}

// Bitset | (or); union of receiver and another set.
func (b *Bitset64) Union(ob *Bitset64) (result *Bitset64) {
	b, ob = sortByLength64(b, ob)
	result = ob.Clone()
	szl := ob.wordCount()
	l := uint64(len(b.b))
	for i := uint64(0); i < l; i++ {
		if i >= szl {
			break
		}
		result.b[i] = b.b[i] | ob.b[i]
	}
	return
}

// Bitset ^ (xor); symmetric difference of receiver and another set.
func (b *Bitset64) SymmetricDifference(ob *Bitset64) (result *Bitset64) {
	b, ob = sortByLength64(b, ob)
	// ob is bigger, so clone it
	result = ob.Clone()
	szl := b.wordCount()
	l := uint64(len(b.b))
	for i := uint64(0); i < l; i++ {
		if i >= szl {
			break
		}
		result.b[i] = b.b[i] ^ ob.b[i]
	}
	return
}

// Return true if the bitset's length is a multiple of the word size.
func (b *Bitset64) isEven() bool {
	return (b.n % sw_64) == 0
}

// Clean last word by setting unused bits to 0.
func (b *Bitset64) cleanLastWord() {
	if !b.isEven() {
		b.b[wordsNeeded64(b.n)-1] &= (hff_64 >> (sw_64 - (b.n % sw_64)))
	}
}

// Return the (local) complement of a bitset (up to n bits).
// func (b *Bitset64) Complement() (result *Bitset64) {
// 	result = New64(b.n)
// 	for i, w := range b.b {
// 		result.b[i] = ^(w)
// 	}
// 	result.cleanLastWord()
// 	return
// }

// Returns true if all bits in the bitset are set.
func (b *Bitset64) All() bool {
	return b.Count() == b.n
}

// Returns true if no bit in the bitset is set.
func (b *Bitset64) None() bool {
	for _, w := range b.b {
		if w > 0 {
			return false
		}
	}
	return true
}

// Return true if any bit in the bitset is set.
func (b *Bitset64) Any() bool {
	return !b.None()
}

// Get a string representation of the words in the bitset.
func (b *Bitset64) String() string {
	f := bytes.NewBufferString("")
	for i := int(wordsNeeded64(b.n) - 1); i >= 0; i-- {
		fmt.Fprintf(f, "%064b.", b.b[i])
	}
	return f.String()
}

// Make a new bitset with a starting capacity of n bits. The bitset expands
// automatically.
func New64(n uint64) *Bitset64 {
	nWords := wordsNeeded64(n)
	if nWords > math.MaxInt32-1 {
		panic(fmt.Sprintf("Bitset64 needs %d %d-bit words to store %d bits, but slices cannot hold more than %d items.", nWords, sw_64, n, math.MaxInt32-1))
	}
	b := &Bitset64{
		n,
		make([]uint64, nWords),
	}
	return b
}
