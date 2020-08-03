package topbuffer

// TopBuffer store sortered according with fgt function (comparator) slice
type TopBuffer struct {
	buffer []int
	size   int
	fgt    func(lhs, rhs int) bool
}

// Add try insert element v in buffer and returns position of it element
// if inserted and -1 otherwise
func (b *TopBuffer) Add(v int) int {
	insertIndex := -1
	// if buffer has vacant places append new value
	if len(b.buffer) < b.size {
		b.buffer = append(b.buffer, v)
		insertIndex = len(b.buffer) - 1
	} else {
		// if new element large then min, overwrite min element
		if b.fgt(v, b.buffer[len(b.buffer) - 1]) {
			insertIndex = len(b.buffer) - 1
			b.buffer[insertIndex] = v
		}
	}

	// swap all elements with new till we goes to the 0-place of find larger element
	if insertIndex > 0 {
		for insertIndex > 0 && b.fgt(b.buffer[insertIndex], b.buffer[insertIndex-1]) {
			b.buffer[insertIndex-1], b.buffer[insertIndex] = b.buffer[insertIndex], b.buffer[insertIndex-1]
			insertIndex--
		}
	}

	return insertIndex
}

// Get returns  sortered slice
func (b TopBuffer) Get() []int {
	return b.buffer
}
