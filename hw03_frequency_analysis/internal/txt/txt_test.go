package txt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextChunk(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		input    string
		delim    string
		maxLen   int
		expected string
		err      error
	}{
		{
			name:     `empty`,
			input:    ``,
			delim:    ` `,
			maxLen:   10,
			expected: ``,
		},
		{
			//         123456789012345678901234567890
			name:     `one word`,
			input:    `one`,
			delim:    ` `,
			maxLen:   10,
			expected: `one`,
		},
		{
			//         123456789012345678901234567890
			name:     `two words in single chunk`,
			input:    `one two`,
			delim:    ` `,
			maxLen:   10,
			expected: `one two`,
		},
		{
			//         123456789012345678901234567890
			name:     `three words in two chunks`,
			input:    `one two three`,
			delim:    ` `,
			maxLen:   10,
			expected: `one two`,
		},
		{
			//         123456789012345678901234567890
			name:     `large then chunk sequense of symbols`,
			input:    `onetwothree`,
			delim:    ` `,
			maxLen:   10,
			expected: ``,
			err:      fmt.Errorf("there is no any delimiter in sequens large, then chunk"),
		},
		{
			//         123456789012345678901234567890
			name:     `chunk starts with space`,
			input:    ` one two three`,
			delim:    ` `,
			maxLen:   10,
			expected: ` one two`,
		},
		{
			//         123456789012345678901234567890
			name:     `chunk ends with space`,
			input:    `one two three `,
			delim:    ` `,
			maxLen:   10,
			expected: `one two`,
		},
		{
			//         123456789012345678901234567890
			name:     `input is only spaces`,
			input:    `               `,
			delim:    ` `,
			maxLen:   10,
			expected: `         `,
		},
		{
			//         123456789012345678901234567890
			name:     `another delimiter`,
			input:    `one_two_three`,
			delim:    `_`,
			maxLen:   10,
			expected: `one_two`,
		},
		{
			//         123456789012345678901234567890
			name:     `several symbols delimiter`,
			input:    `onego twogo three`,
			delim:    `go`,
			maxLen:   15,
			expected: `onego two`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			chr := Chunker{
				buffer: tt.input,
				delim:  tt.delim,
				maxLen: tt.maxLen,
			}
			result, err := chr.NextChunk()
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestChunks(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		input    string
		delim    string
		maxLen   int
		expected []string
		err      error
	}{
		{
			name:     `empty`,
			input:    ``,
			delim:    ` `,
			maxLen:   10,
			expected: []string{},
		},
		{
			name:     `one word`,
			input:    `one`,
			delim:    ` `,
			maxLen:   10,
			expected: []string{`one`},
		},
		{
			name:     `two words in single chunk`,
			input:    `one two`,
			delim:    ` `,
			maxLen:   10,
			expected: []string{`one two`},
		},
		{
			name:     `three words in two chunks`,
			input:    `one two three`,
			delim:    ` `,
			maxLen:   10,
			expected: []string{`one two`, `three`},
		},
		{
			name:     `large then chunk sequense of symbols`,
			input:    `onetwothree`,
			delim:    ` `,
			maxLen:   10,
			expected: []string{},
			err:      fmt.Errorf("there is no any delimiter in sequens large, then chunk"),
		},
		{
			name:     `second chunk is large then maxLen`,
			input:    `one twothreefour`,
			delim:    ` `,
			maxLen:   10,
			expected: []string{`one`},
			err:      fmt.Errorf("there is no any delimiter in sequens large, then chunk"),
		},
		{
			name:     `chunk starts with space`,
			input:    ` one two three`,
			delim:    ` `,
			maxLen:   10,
			expected: []string{` one two`, `three`},
		},
		{
			name:     `chunk ends with space`,
			input:    `one two three `,
			delim:    ` `,
			maxLen:   10,
			expected: []string{`one two`, `three `},
		},
		{
			name:     `input is only spaces`,
			input:    `               `,
			delim:    ` `,
			maxLen:   10,
			expected: []string{`         `, `     `},
		},
		{
			name:     `another delimiter`,
			input:    `one_two_three`,
			delim:    `_`,
			maxLen:   10,
			expected: []string{`one_two`, `three`},
		},
		{
			name:     `left side delimiter splitting`,
			input:    `one_two__three`,
			delim:    `_`,
			maxLen:   10,
			expected: []string{`one_two_`, `three`},
		},
		{
			name:     `several symbols delimiter`,
			input:    `onego twogo three`,
			delim:    `go`,
			maxLen:   15,
			expected: []string{`onego two`, ` three`},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			chr := Chunker{
				buffer: tt.input,
				delim:  tt.delim,
				maxLen: tt.maxLen,
			}
			result, err := Chunks(chr)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}

}
