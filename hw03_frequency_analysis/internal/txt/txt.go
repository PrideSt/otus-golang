package txt

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

var _ io.ByteWriter = (*strings.Builder)(nil)

func writeLastRank(i int, alplabet []byte, b io.ByteWriter) {
	if i > 0 {
		writeLastRank(i/len(alplabet), alplabet, b)
		_ = b.WriteByte(alplabet[i%len(alplabet)])
	}
}

// IntToWord converts given numeric to numeric system with base = len(alphabet) and
// use given letters in alphabet as rank values.
func IntToWord(i int, alplabet []byte) string {
	if i == 0 {
		return string(alplabet[0])
	}
	// var builder io.ByteWriter
	builder := strings.Builder{}

	writeLastRank(i, alplabet, &builder)
	return builder.String()
}

// GenDict returns slice of words with len wordsCnt starts from offset with given step,
// using alplabet for word characters.
func GenDict(wordsCnt int, alplabet []byte, offset int, step int) ([]string, error) {
	result := make([]string, wordsCnt)

	if wordsCnt < 0 {
		return nil, fmt.Errorf("invalid argument: wordsCnt, given negative value %d", wordsCnt)
	}

	if len(alplabet) == 0 {
		return nil, errors.New("invalid argument: alplabet, is empty")
	}

	if offset < 0 {
		return nil, fmt.Errorf("invalid argument: offset, given negative value %d", offset)
	}

	// @todo add test and realisation for negative step
	if step < 0 {
		return nil, fmt.Errorf("invalid argument: step, given negative value %d", step)
	}

	for i := 0; i < wordsCnt; i++ {
		result[i] = IntToWord(offset, alplabet)
		// @todo check offset overlap
		offset += step
	}

	return result, nil
}

// GenText returns string with wordsCnt words in it, takes words from slice of string,
// this function using standard random generator with normal distribution.
func GenText(dict []string, wordsCnt int, r *rand.Rand) string {
	var builder strings.Builder

	for i := 0; i < wordsCnt-1; i++ {
		n := r.Intn(len(dict))
		builder.WriteString(dict[n])
		builder.WriteByte(byte(' '))
	}

	n := r.Intn(len(dict))
	builder.WriteString(dict[n])

	return builder.String()
}
