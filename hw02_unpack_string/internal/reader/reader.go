package reader

import (
	"unicode/utf8"
)

// Reader is wrapper for utf8 decode function.
type Reader struct {
	buffer []byte
	offset int
}

// New returns new Reader instance, given buffer copied.
func Make(buffer []byte, offset int) Reader {
	ownBuffer := make([]byte, len(buffer))
	copy(ownBuffer, buffer)

	return Reader{ownBuffer, offset}
}

// GetNext returns next rune from internal buffer.
func (r *Reader) GetNext() (rune, int) {
	rn, sz := utf8.DecodeRune(r.buffer[r.offset:])
	r.offset += sz

	return rn, r.offset
}

// IsEOF returns false when offset goes till end of internal buffer.
func (r *Reader) IsEOF() bool {
	return r.offset >= len(r.buffer)
}
