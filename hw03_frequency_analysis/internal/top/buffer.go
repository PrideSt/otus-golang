package top

import (
	"fmt"
	"sort"
)

type FreqEntry struct {
	Word string
	Freq int
}

// Buffer store sortered according with fgt function (comparator) slice.
// Using Buffer take much less memory and cpu for taking reault, see benchmark results below.
type Buffer struct {
	len    int
	buffer []FreqEntry
	less   func(lhs, rhs FreqEntry) bool
}

type Interface interface {
	sort.Interface
	Add(v FreqEntry) int
	Get() []FreqEntry
}

// Make Buffer with len = size and element comparator fgt.
func NewBuffer(size int, less func(lhs, rhs FreqEntry) bool) *Buffer {
	buffer := make([]FreqEntry, size)
	instance := Buffer{
		buffer: buffer,
		less:   less,
	}

	return &instance
}

// Len is the number of elements in the collection, implementation of sort.Len().
func (b *Buffer) Len() int {
	return len(b.buffer)
}

// Less compare function, implementation of sort.Less().
func (b *Buffer) Less(i, j int) bool {
	if i > b.len || j > b.len {
		panic(fmt.Sprintf("Less func get index [i,j]: [%d, %d] out of range, buffer len: %d, buffer cap: %d", i, j, b.len, len(b.buffer)))
	}

	return b.less(b.buffer[i], b.buffer[j])
}

// Swap swaps the elements with indexes i and j.
func (b *Buffer) Swap(i, j int) {
	b.buffer[i], b.buffer[j] = b.buffer[j], b.buffer[i]
}

// Add try insert element v in buffer and returns position of it element
// if inserted and -1 otherwise.
func (b *Buffer) Add(v FreqEntry) int {
	insertIndex := -1
	// if buffer has vacant places append new value
	if b.len < len(b.buffer) {
		b.buffer[b.len] = v
		insertIndex = b.len
		b.len++
	} else if b.less(v, b.buffer[b.len-1]) {
		// if new element large then min, overwrite min element
		insertIndex = b.len - 1
		b.buffer[insertIndex] = v
	}

	// swap all elements with new till we goes to the 0-place of find larger element
	if insertIndex > 0 {
		for insertIndex > 0 && b.Less(insertIndex, insertIndex-1) {
			b.Swap(insertIndex-1, insertIndex)
			insertIndex--
		}
	}

	return insertIndex
}

// FindTopN pass every antry from ferq dictionary (map of word -> count in text) through Interface- ordered buffern with
// topResultsCnt elements. Returns ordered slice of pair {word, freq}.
func FindTopN(dict map[string]int, topResultsCnt int, chTerminate <-chan struct{}) []FreqEntry {
	buffer := NewBuffer(topResultsCnt, func(lhs, rhs FreqEntry) bool {
		return lhs.Freq > rhs.Freq
	})

	i := 0
	for w, fr := range dict {
		if i%100 == 0 {
			// we should terminate while finding top, check it every 100 iterations
			select {
			case <-chTerminate:

				return []FreqEntry{}
			default:
			}
		}
		buffer.Add(FreqEntry{Word: w, Freq: fr})
		i++
	}

	return buffer.Get()
}

// Get returns  sortered slice.
func (b *Buffer) Get() []FreqEntry {
	return b.buffer[:b.len]
}

// sort.Interface need for brnchmark.
type FreqSlice []FreqEntry

func (p FreqSlice) Len() int           { return len(p) }
func (p FreqSlice) Less(i, j int) bool { return p[i].Freq < p[j].Freq }
func (p FreqSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
