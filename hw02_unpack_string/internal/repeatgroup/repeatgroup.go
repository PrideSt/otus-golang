package repeatgroup

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Parser can create Unpacker from string, create internal view from encoded string
type Parser interface {
	ParseString(str string) (Unpacker, error)
}

// Unpacker can create string from Unpacker, assembly string from internal view
type Unpacker interface {
	Unpack() (string, error)
}

// RepeatGroup basic struct of internal representation
type repeatGroup struct {
	buffer []byte
	repeatCnt int
}

// GroupStorage contain all groups
type GroupStorage struct {
	rgs []repeatGroup
}

// AddRepeatGroup new entry to storage
func (gs *GroupStorage) addRepeatGroup(newRG repeatGroup) (newSize int){
	gs.rgs = append(gs.rgs, newRG)
	return len(gs.rgs)
}

// Add new RepeatGroup entry to storage created it from butes (copy)
func (gs *GroupStorage) add(b []byte, times int) (newSize int){
	chunk := make([]byte, len(b))
	copy(chunk, b)
	gs.addRepeatGroup(repeatGroup{chunk, times})
	return len(gs.rgs)
}

// as unicode.Digit, but allow only ascii digits [0-9]
func isASCIIDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isEscapeSymbol(r rune) bool {
	return r == '\\'
}

// @see https://www.fileformat.info/info/unicode/category/Sk/list.htm
func isSymbolModifier(r rune) bool  {
	return unicode.In(r, unicode.Sk)
}

// @see https://www.fileformat.info/info/unicode/category/Mn/list.htm
func isMarkNonspacing(r rune) bool {
	return unicode.In(r, unicode.Mn)
}

// @see https://github.com/golang/go/blob/master/src/unicode/tables.go#L5952
// \u200C Zero Width Non-Joiner
// \u200D Zero Width Joiner
func isJoiner(r rune) bool {
	return  unicode.In(r, unicode.Join_Control)
}

// ParseString converts input string to internal state
func ParseString(input string) (Unpacker, error) {
	var gs GroupStorage
	// store chunk
	var buffer bytes.Buffer

	offset := 0

	for offset < len(input) {
		r, runeSize := utf8.DecodeRuneInString(input[offset:])
		offset += runeSize

		if isASCIIDigit(r) {
			if buffer.Len() == 0 {
				return GroupStorage{}, fmt.Errorf("invalid format %s, offset: %d, digit can't be the first symbol", input, offset)
			}

			gs.add(buffer.Bytes(), int(r - '0'))
			buffer.Reset()
		} else {
			// put all symbol moifiers to buffer
			if isSymbolModifier(r) || isMarkNonspacing(r) {
				buffer.WriteRune(r)
				continue
			}

			if isJoiner(r) {
				// write joiner
				buffer.WriteRune(r)

				// write next symbol
				if offset < len(input) {
					// overwrite variables in outer scope
					r, runeSize = utf8.DecodeRuneInString(input[offset:])
					offset += runeSize
				} else {
					return GroupStorage{}, fmt.Errorf("invalid format %s, offset: %d, input string can't ends with joiner", input, offset)
				}

				buffer.WriteRune(r)
				continue
			}

			// flush dry buffer
			if buffer.Len() > 0 {
				gs.add(buffer.Bytes(), 1)
				buffer.Reset()
			}

			// if escape symbol given, repeat reading
			if isEscapeSymbol(r) {
				// if escape is last symbol do nothing, add them to buffer like any another
				// otherwise read next symbol
				if offset < len(input) {
					// overwrite variables in outer scope
					r, runeSize = utf8.DecodeRuneInString(input[offset:])
					offset += runeSize
				}
			}
			buffer.WriteRune(r)
		}
	}

	if buffer.Len() > 0 {
		gs.add(buffer.Bytes(), 1)
	}

	return gs, nil
}

// Unpack converts internal state to string
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