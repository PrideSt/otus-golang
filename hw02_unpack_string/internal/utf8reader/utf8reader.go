package utf8reader

import (
	"unicode/utf8"
)

// Reader is wrapper for utf8 decode function.
type Reader struct {
	buffer []byte
	offset int
}

// Make create new instance of Reader.
func Make(b []byte) Reader {
	tmp := Reader{buffer: b}
	return tmp
}

// GetNext returns next rune from internal buffer.
func (r *Reader) GetNext() (rune, int) {
	rn, sz := utf8.DecodeRune(r.buffer[r.offset:])
	r.offset += sz

	return rn, r.offset
}

// IsNotEOF returns false when offset goes till end of internal buffer.
func (r *Reader) IsNotEOF() bool {
	return r.offset < len(r.buffer)
}
