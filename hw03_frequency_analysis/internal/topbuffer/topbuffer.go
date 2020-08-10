package topbuffer

import (
	"fmt"
	"sort"
)

type FreqEntry struct {
	Word string
	Freq int
}

// TopBuffer store sortered according with fgt function (comparator) slice.
// Using TopBuffer take much less memory and cpu for taking reault, see benchmark results below.
type TopBuffer struct {
	len    int
	buffer []FreqEntry
	less   func(lhs, rhs FreqEntry) bool
}

type TopBufferInterface interface {
	sort.Interface
	Add(v FreqEntry) int
	Get() []FreqEntry
}

// Make TopBuffer with len = size and element comparator fgt.
func New(size int, less func(lhs, rhs FreqEntry) bool) *TopBuffer {
	buffer := make([]FreqEntry, size)
	// buffer := []FreqEntry{}
	instance := TopBuffer{
		buffer: buffer,
		less:   less,
	}

	return &instance
}

// Len is the number of elements in the collection, implementation of sort.Len().
func (b *TopBuffer) Len() int {
	return len(b.buffer)
}

// Less compare function, implementation of sort.Less().
func (b *TopBuffer) Less(i, j int) bool {
	if i > b.len || j > b.len {
		panic(fmt.Sprintf("Less func get index [i,j]: [%d, %d] out of range, buffer len: %d, buffer cap: %d", i, j, b.len, len(b.buffer)))
	}
	return b.less(b.buffer[i], b.buffer[j])
}

// Swap swaps the elements with indexes i and j.
func (b *TopBuffer) Swap(i, j int) {
	b.buffer[i], b.buffer[j] = b.buffer[j], b.buffer[i]
}

// Add try insert element v in buffer and returns position of it element
// if inserted and -1 otherwise.
func (b *TopBuffer) Add(v FreqEntry) int {
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

// Get returns  sortered slice.
func (b *TopBuffer) Get() []FreqEntry {
	return b.buffer[:b.len]
}
