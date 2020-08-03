package repeatgroup

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/PrideSt/otus-golang/hw02_unpack_string/internal/reader"
)

// Parser can create Unpacker from string, create internal view from encoded string.
type Parser interface {
	ParseString(str string) (Unpacker, error)
}

// RepeatGroup basic struct of internal representation.
type repeatGroup struct {
	buffer    []byte
	repeatCnt int
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

// ParseString converts input string to internal state.
func ParseString(input string) (Unpacker, error) {
	var gs GroupStorage
	var buffer bytes.Buffer
	reader := reader.Make([]byte(input), 0)

	for !reader.IsEOF() {
		r, offset := reader.GetNext()

		switch {
		case isASCIIDigit(r):
			if buffer.Len() == 0 {
				return GroupStorage{}, fmt.Errorf("invalid format %s, offset: %d, digit can't be the first symbol", input, offset)
			}

			gs.flushBuffer(&buffer, int(r-'0'))

			continue
		case isEscapeSymbol(r):
			// if escape is last symbol do nothing, add them to buffer like any another
			// otherwise read next symbol
			gs.flushBuffer(&buffer, 1)
			if !reader.IsEOF() {
				// overwrite variables in outer scope
				r, _ = reader.GetNext()
			}
		case isJoiner(r):
			// write joiner
			buffer.WriteRune(r)

			// write next symbol
			if !reader.IsEOF() {
				// overwrite variables in outer scope
				r, _ = reader.GetNext()
			} else {
				return GroupStorage{}, fmt.Errorf("invalid format %s, offset: %d, input string can't ends with joiner", input, offset)
			}
		case isSymbolModifier(r) || isMarkNonspacing(r):
		default:
			gs.flushBuffer(&buffer, 1)
		}
		buffer.WriteRune(r)
	}

	gs.flushBuffer(&buffer, 1)

	return gs, nil
}
