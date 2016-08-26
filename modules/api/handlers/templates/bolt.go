package templates

// This file contains low-level bolt structs that are used for reading from
// bolt database files.

import (
	"fmt"
	"unsafe"

	"github.com/boltdb/bolt"
)

const pageHeaderSize = int(unsafe.Offsetof(((*page)(nil)).ptr))
const branchPageElementSize = int(unsafe.Sizeof(branchPageElement{}))
const leafPageElementSize = int(unsafe.Sizeof(leafPageElement{}))

const maxAllocSize = 0xFFFFFFF
const maxNodesPerPage = 65535

const (
	branchPageFlag   = 0x01
	leafPageFlag     = 0x02
	metaPageFlag     = 0x04
	freelistPageFlag = 0x10
)

const (
	bucketLeafFlag = 0x01
)

type pgid uint64
type txid uint64

type page struct {
	id       pgid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}

type stats struct {
	inuse       int
	alloc       int
	utilization float64
	histogram   map[int]int
}

// typ returns a human readable page type string used for debugging.
func (p *page) typ() string {
	if (p.flags & branchPageFlag) != 0 {
		return "branch"
	} else if (p.flags & leafPageFlag) != 0 {
		return "leaf"
	} else if (p.flags & metaPageFlag) != 0 {
		return "meta"
	} else if (p.flags & freelistPageFlag) != 0 {
		return "freelist"
	}
	return fmt.Sprintf("unknown<%02x>", p.flags)
}

func (p *page) meta() *meta {
	return (*meta)(unsafe.Pointer(&p.ptr))
}

func (p *page) leafPageElement(index uint16) *leafPageElement {
	n := &((*[maxNodesPerPage]leafPageElement)(unsafe.Pointer(&p.ptr)))[index]
	return n
}

func (p *page) branchPageElement(index uint16) *branchPageElement {
	return &((*[maxNodesPerPage]branchPageElement)(unsafe.Pointer(&p.ptr)))[index]
}

// stats calcuates statistics for a page.
func (p *page) stats(pageSize int) stats {
	var s stats
	s.alloc = (int(p.overflow) + 1) * pageSize
	s.inuse = p.inuse()

	// Calculate space utilitization
	if s.alloc > 0 {
		s.utilization = float64(s.inuse) / float64(s.alloc)
	}

	return s
}

// inuse returns the number of bytes used in a given page.
func (p *page) inuse() int {
	var n int
	if (p.flags & leafPageFlag) != 0 {
		n = pageHeaderSize
		if p.count > 0 {
			n += leafPageElementSize * int(p.count-1)
			e := p.leafPageElement(p.count - 1)
			n += int(e.pos + e.ksize + e.vsize)
		}

	} else if (p.flags & branchPageFlag) != 0 {
		n = pageHeaderSize
		if p.count > 0 {
			n += branchPageElementSize * int(p.count-1)
			e := p.branchPageElement(p.count - 1)
			n += int(e.pos + e.ksize)
		}
	}
	return n
}

// usage calculates a histogram of page sizes within nested pages.
func usage(tx *bolt.Tx, pgid pgid) map[int]int {
	m := make(map[int]int)
	forEachPage(tx, pgid, func(p *page) {
		m[p.inuse()]++
	})
	return m
}

// branchPageElement represents a node on a branch page.
type branchPageElement struct {
	pos   uint32
	ksize uint32
	pgid  pgid
}

// key returns a byte slice of the node key.
func (n *branchPageElement) key() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return buf[n.pos : n.pos+n.ksize]
}

// leafPageElement represents a node on a leaf page.
type leafPageElement struct {
	flags uint32
	pos   uint32
	ksize uint32
	vsize uint32
}

// key returns a byte slice of the node key.
func (n *leafPageElement) key() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return buf[n.pos : n.pos+n.ksize]
}

// value returns a byte slice of the node value.
func (n *leafPageElement) value() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return buf[n.pos+n.ksize : n.pos+n.ksize+n.vsize]
}

type meta struct {
	magic    uint32
	version  uint32
	pageSize uint32
	flags    uint32
	root     bucket
	freelist pgid
	pgid     pgid
	txid     txid
	checksum uint64
}

// bucket_ represents the bolt.Bucket type.
type bucket_ struct {
	*bucket
}

type bucket struct {
	root     pgid
	sequence uint64
}

type tx struct {
	writable bool
	managed  bool
	db       uintptr
	meta     *meta
	root     bucket
	// remaining fields not used.
}

// find locates a page using either a set of page indices or a direct page id.
// It returns a list of parent page numbers if the indices are used.
// It also returns the page reference.
func find(tx *bolt.Tx, directID int, indexes []int) (*page, []pgid, error) {
	// If a direct ID is provided then just use it.
	if directID != 0 {
		return pageAt(tx, pgid(directID)), nil, nil
	}

	// Otherwise traverse the pages index.
	ids, err := pgids(tx, indexes)
	if err != nil {
		return nil, nil, err
	}

	return pageAt(tx, ids[len(ids)-1]), ids, nil
}

// retrieves the page from a given transaction.
func pageAt(tx *bolt.Tx, id pgid) *page {
	info := tx.DB().Info()
	return (*page)(unsafe.Pointer(info.Data + uintptr(info.PageSize*int(id))))
}

// forEachPage recursively iterates over all pages starting at a given page.
func forEachPage(tx *bolt.Tx, pgid pgid, fn func(*page)) {
	p := pageAt(tx, pgid)
	fn(p)

	if (p.flags & leafPageFlag) != 0 {
		for i := 0; i < int(p.count); i++ {
			if e := p.leafPageElement(uint16(i)); (e.flags & bucketLeafFlag) != 0 {
				if b := (*bucket)(unsafe.Pointer(&e.value()[0])); b.root != 0 {
					forEachPage(tx, b.root, fn)
				}
			}
		}
	} else if (p.flags & branchPageFlag) != 0 {
		for i := 0; i < int(p.count); i++ {
			forEachPage(tx, p.branchPageElement(uint16(i)).pgid, fn)
		}
	}
}
