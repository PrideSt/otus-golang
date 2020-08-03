package txt

import (
	"fmt"
	"strings"
)

// Chunker store string buffer inside and returns parts of this buffer
// splited by delim symbol/symbols by parts less or equal then size then
// maxLen, delimiter not appends to chunks
type Chunker struct {
	buffer string
	delim  string
	maxLen int
}

// Chunks return slice of buffer aparts
func Chunks(c Chunker) ([]string, error) {
	result := []string{}
	for {
		chunk, err := c.NextChunk()
		if err != nil {
			return result, err
		}
		if chunk == "" {
			break
		}
		result = append(result, chunk)
	}
	return result, nil
}

// NextChunk return first max avaliable chunk size less of equal size then c.maxLen
// splited by c.delim and chop c.buffer to new position. Delimiter removes from chunks,
// but another delimiters on chunk boundary not
// returns zero length string and error when meet byte sequense large then c.maxLen
// without any delimiters
func (c *Chunker) NextChunk() (string, error) {
	var chunk string
	if len(c.buffer) <= c.maxLen {
		chunk = c.buffer
		c.buffer = ""
	} else {
		newOffset := strings.LastIndex(c.buffer[:c.maxLen], c.delim)

		if newOffset < 0 {
			return "", fmt.Errorf("there is no any delimiter in sequens large, then chunk")
		}

		chunk = c.buffer[:newOffset]
		if (newOffset + len(c.delim) < len(c.buffer)) {
			c.buffer = c.buffer[newOffset + len(c.delim):]
		} else {
			c.buffer = ""
		}
	}

	return chunk, nil
}
