package repeatgroup

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/PrideSt/otus-golang/hw02_unpack_string/internal/utf8reader"
)

// Parser can create Unpacker from string, create internal view from encoded string.
type Parser interface {
	ParseString(str string) (Unpacker, error)
}

// Unpacker can create string from Unpacker, assembly string from internal view.
type Unpacker interface {
	Unpack() (string, error)
}

// RepeatGroup basic struct of internal representation.
type repeatGroup struct {
	buffer    []byte
	repeatCnt int
}

// GroupStorage contain all groups.
type GroupStorage struct {
	rgs []repeatGroup
}

// AddRepeatGroup new entry to storage.
func (gs *GroupStorage) addRepeatGroup(newRG repeatGroup) (newSize int) {
	gs.rgs = append(gs.rgs, newRG)
	return len(gs.rgs)
}

// Add new RepeatGroup entry to storage created it from butes (copy).
func (gs *GroupStorage) add(b []byte, times int) (newSize int) {
	chunk := make([]byte, len(b))
	copy(chunk, b)
	gs.addRepeatGroup(repeatGroup{chunk, times})
	return len(gs.rgs)
}

// as unicode.Digit, but allow only ascii digits [0-9].
func isASCIIDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isEscapeSymbol(r rune) bool {
	return r == '\\'
}

// @see https://www.fileformat.info/info/unicode/category/Sk/list.htm
func isSymbolModifier(r rune) bool {
	return unicode.In(r, unicode.Sk)
}

// @see https://www.fileformat.info/info/unicode/category/Mn/list.htm
func isMarkNonspacing(r rune) bool {
	return unicode.In(r, unicode.Mn)
}

// isJoiner return true in case when rune is Join code-point.
func isJoiner(r rune) bool {
	// @see https://github.com/golang/go/blob/master/src/unicode/tables.go#L5952
	// \u200C Zero Width Non-Joiner
	// \u200D Zero Width Joiner
	return unicode.In(r, unicode.Join_Control)
}

// flushBuffer create repeatGroup from buffer and cnt and flush buffer.
func flushBuffer(gs *GroupStorage, buffer *bytes.Buffer, cnt int) {
	if buffer.Len() > 0 {
		gs.add(buffer.Bytes(), cnt)
		buffer.Reset()
	}
}

// ParseString converts input string to internal state.
func ParseString(input string) (Unpacker, error) {
	var gs GroupStorage
	var buffer bytes.Buffer
	reader := utf8reader.Make([]byte(input))

	for reader.IsNotEOF() {
		r, offset := reader.GetNext()

		switch {
		case isASCIIDigit(r):
			if buffer.Len() == 0 {
				return GroupStorage{}, fmt.Errorf("invalid format %s, offset: %d, digit can't be the first symbol", input, offset)
			}

			flushBuffer(&gs, &buffer, int(r-'0'))
			continue
		case isEscapeSymbol(r):
			// if escape is last symbol do nothing, add them to buffer like any another
			// otherwise read next symbol
			flushBuffer(&gs, &buffer, 1)
			if reader.IsNotEOF() {
				// overwrite variables in outer scope
				r, _ = reader.GetNext()
			}
		case isJoiner(r):
			// write joiner
			buffer.WriteRune(r)

			// write next symbol
			if reader.IsNotEOF() {
				// overwrite variables in outer scope
				r, _ = reader.GetNext()
			} else {
				return GroupStorage{}, fmt.Errorf("invalid format %s, offset: %d, input string can't ends with joiner", input, offset)
			}
		case isSymbolModifier(r) || isMarkNonspacing(r):
		default:
			flushBuffer(&gs, &buffer, 1)
		}
		buffer.WriteRune(r)
	}

	flushBuffer(&gs, &buffer, 1)
	return gs, nil
}

// Unpack converts internal state to string.
func (gs GroupStorage) Unpack() (string, error) {
	var result bytes.Buffer
	for _, gr := range gs.rgs {
		if gr.repeatCnt > 0 {
			newChunk := bytes.Repeat(gr.buffer, gr.repeatCnt)
			writedLen, err := result.Write(newChunk)
			if err != nil {
				return "", fmt.Errorf("unable to write into buffer, %w", err)
			}
			if writedLen < len(gr.buffer) {
				return "", fmt.Errorf("buffer writes partially")
			}
		}
	}

	return result.String(), nil
}
