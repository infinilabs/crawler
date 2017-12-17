package ui

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/boltdb/bolt"
)

// tostr converts a byte slice to a string if all characters are printable.
// otherwise prints the hex representation.
func tostr(b []byte) string {
	// Check if byte slice is utf-8 encoded.
	if !utf8.Valid(b) {
		return fmt.Sprintf("%x", b)
	}

	// Check every rune to see if it's printable.
	var s = string(b)
	for _, ch := range s {
		if !unicode.IsPrint(ch) {
			return fmt.Sprintf("%x", b)
		}
	}

	return s
}

func trunc(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

// traverses the page hierarchy by index and returns associated page ids.
// returns an error if an index is out of range.
func pgids(t *bolt.Tx, indexes []int) ([]pgid, error) {
	tx := (*tx)(unsafe.Pointer(t))

	p := pageAt(t, tx.meta.root.root)
	ids := []pgid{tx.meta.root.root}
	for _, index := range indexes[1:] {
		if uint16(index) >= p.count {
			return nil, fmt.Errorf("out of range")
		}

		if (p.flags & branchPageFlag) != 0 {
			e := p.branchPageElement(uint16(index))
			ids = append(ids, e.pgid)
			p = pageAt(t, e.pgid)

		} else if (p.flags & leafPageFlag) != 0 {
			// Only non-inline buckets on leaf pages can be traversed.
			e := p.leafPageElement(uint16(index))
			if (e.flags & bucketLeafFlag) == 0 {
				return nil, fmt.Errorf("index not a bucket")
			}

			b := (*bucket)(unsafe.Pointer(&e.value()[0]))
			if (e.flags & bucketLeafFlag) == 0 {
				return nil, fmt.Errorf("index is an inline bucket")
			}

			ids = append(ids, b.root)
			p = pageAt(t, b.root)
		} else {
			return nil, fmt.Errorf("invalid page type: %s" + p.typ())
		}
	}
	return ids, nil
}

func pagelink(indexes []int) string {
	var a []string
	for _, index := range indexes[1:] {
		a = append(a, strconv.Itoa(index))
	}
	return "boltdb?index=" + strings.Join(a, ":")
}

func subpagelink(indexes []int, index int) string {
	var tmp = make([]int, len(indexes))
	copy(tmp, indexes)
	tmp = append(tmp, index)
	return pagelink(tmp)
}

// borrowed from: https://github.com/dustin/go-humanize
func comma(v int) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(int64(v)%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

// maxHistogramN is the maximum number of buckets that can be used for the histogram.
const maxHistogramN = 100

// bucketize converts a map of observations into a histogram with a
// smaller set of buckets using the square root choice method.
func bucketize(m map[int]int) (mins, maxs, values []int) {
	if len(m) == 0 {
		return nil, nil, nil
	} else if len(m) == 1 {
		for k, v := range m {
			return []int{k}, []int{k}, []int{v}
		}
	}

	// Retrieve sorted set of keys.
	var keys []int
	var vsum int
	for k, v := range m {
		keys = append(keys, k)
		vsum += v
	}
	sort.Ints(keys)

	// Determine min/max for 5-95 percentile.
	var pmin, pmax, x int
	for _, k := range keys {
		v := m[k]

		// Grab the min when we cross the 5% threshold and 95% threshold.
		if (x*100)/vsum < 5 && ((x+v)*100)/vsum >= 5 {
			pmin = k
		}
		if (x*100)/vsum < 95 && ((x+v)*100)/vsum >= 95 {
			pmax = k
		}

		x += m[k]
	}
	min, max := keys[0], keys[len(keys)-1]

	// Calculate number of buckets and step size.
	n := int(math.Ceil(math.Sqrt(float64(vsum))))
	if n > maxHistogramN {
		n = maxHistogramN
	}
	step := float64(pmax-pmin) / float64(n)

	// Bucket everything.
	for i := 0; i < n; i++ {
		var kmin, kmax int
		if i == 0 {
			kmin = min
		} else {
			kmin = int(math.Floor(float64(pmin) + step*float64(i)))
		}
		if i == n-1 {
			kmax = max
		} else {
			kmax = int(math.Floor(float64(pmin)+step*float64(i+1))) - 1
		}
		value := 0

		for k, v := range m {
			if (i == 0 || k >= kmin) && (i == (n-1) || k <= kmax) {
				value += v
			}
		}

		mins = append(mins, kmin)
		maxs = append(maxs, kmax)
		values = append(values, value)
	}

	return
}
