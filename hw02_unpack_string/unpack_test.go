package hw02_unpack_string //nolint:golint,stylecheck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PrideSt/otus-golang/hw02_unpack_string/internal/repeatgroup"
)

func TestUnpack(t *testing.T) {
	for _, tst := range [...]struct{
		name	 string
		input    string
		expected string
		err      error
	}{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      fmt.Errorf("invalid format 3abc, offset: 1, digit can't be the first symbol"),
		},
		{
			input:    "45",
			expected: "",
			err:      fmt.Errorf("invalid format 45, offset: 1, digit can't be the first symbol"),
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      fmt.Errorf("invalid format aaa10b, offset: 5, digit can't be the first symbol"),
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
	} {
		gs, err := repeatgroup.ParseString(tst.input)
		require.Equal(t, tst.err, err)

		result, _ := gs.Unpack()
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for _, tst := range [...]struct{
		name	 string
		input    string
		expected string
		err      error
	}{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
	} {
		gs, err := repeatgroup.ParseString(tst.input)
		require.Equal(t, tst.err, err)

		result, _ := gs.Unpack()
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithZeroByte(t *testing.T) {
	// t.Skip()
	for _, tt := range [...]struct{
		name	 string
		input    string
		expected string
		err      error
	} {
		{
			name:     `zero byte first`,
			input:    "\x00a2",
			expected: "\x00aa",
			// err: fmt.Errorf("invalid format \\x00a2, offset: 4, digit can't be the first symbol"),
		},
		{
			name:     `zero byte last`,
			input:    "a2\x00",
			expected: "aa\x00",
		},
		{
			name:     `zero byte repeat`,
			input:    "a2\x002b2",
			expected: "aa\x00\x00bb",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gs, err := repeatgroup.ParseString(tt.input)
			require.Equal(t, tt.err, err)

			result, _ := gs.Unpack()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnpackMyltiButeRune(t *testing.T) {
	for _, tt := range [...]struct{
		name	 string
		input    string
		expected string
		err      error
	}{
		{
			name: 	  "cyrillic string",
			input:    "ы2ю3ф4ъ0",
			expected: "ыыюююфффф",
			err:      nil,
		},
		{
			name: 	  "cyrillic string",
			input:    "ё2Ё3",
			expected: "ёёЁЁЁ",
			err:      nil,
		},
		{
			name: 	  "Hiragana (Japanese) string",
			input:    "か2あ3ぃ4ぅ0",
			expected: "かかあああぃぃぃぃ",
			err:      nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T){
			gs, err := repeatgroup.ParseString(tt.input)
			require.Equal(t, tt.err, err)

			result, _ := gs.Unpack()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnpackMyltiRuneCharacters(t *testing.T) {
	for _, tt := range [...]struct{
		name	 string
		input    string
		expected string
		err      error
	} {
		// tests with unicode
		{
			name: 	  `Single code-point emoji without skin tone modifier`,
			input:    `a2👨2`,
			expected: `aa👨👨`,
		},
		{
			name: 	  `Multy code-point emoji (with skin tone modifier)`,
			input:    `a3👨🏽2`,
			//       👨\xF0\x9F\x91\xA8 // person, code-point 1
			//                       🏽 \xF0\x9F\x8F\xBD // skin tone fitz-5, code-point 2
			expected: `aaa👨🏽👨🏽`,
		},
		{
			name: 	  `Multy code-point emoji (with skin tone modifier and zero width joiner)`,
			input:    `b1👨🏾‍🚀2`,
			// 👨🏾‍🚀 =   "\xF0\x9F\x91\xA8\xF0\x9F\x8F\xBE\xE2\x80\x8B\xF0\x9F\x9A\x80",
			//       👨\xF0\x9F\x91\xA8 // person
			//                        🏽 \xF0\x9F\x8F\xBD // skin tone fitz-5
			//                                          \xE2\x80\x8B // zero width joiner
			//                                                    🚀\xF0\x9F\x9A\x80 //rocket
			expected: `b👨🏾‍🚀👨🏾‍🚀`,
		},
		{
			name: 	  `Multy code-point emoji with several modifiers`,
			input:    "a1e\u0301\u03012",
			expected: `aé́é́`,
		},
		{
			name: 	  `Multy code-point with letter modifier`,
			input:    "a1e\u02EF\u02EF2",
			expected: `ae˯˯e˯˯`,
		},
		{
			name: `Non ascii digits in repeat count (MATHEMATICAL MONOSPACE DIGIT THREE)`,
			input: `a𝟹`,
			// input: "a\xF0\x9D\x9F\xB9",
			// it looks like a3 (@see https://www.fileformat.info/info/unicode/char/1d7f9/index.htm)
			expected: `a𝟹`,
		},
		{
			name: `Repeat non ascii digit (MATHEMATICAL MONOSPACE DIGIT THREE)`,
			input: `a𝟹2`,
			// input: "a\xF0\x9D\x9F\xB92",
			// it looks like a3 (@see https://www.fileformat.info/info/unicode/char/1d7f9/index.htm)
			expected: `a𝟹𝟹`,
		},
	} {
		t.Run(tt.name, func(t *testing.T){
			gs, err := repeatgroup.ParseString(tt.input)
			require.Equal(t, tt.err, err)

			result, err := gs.Unpack()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
