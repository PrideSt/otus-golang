package repeatgroup

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		input    string
		expected GroupStorage
		err      error
	}{
		// trivial tests
		{
			name:     `Empty string`,
			input:    ``,
			expected: GroupStorage{},
		},
		{
			name:     `Simple char`,
			input:    `a`,
			expected: GroupStorage{[]repeatGroup{{[]byte("a"), 1}}},
		},
		{
			name:  `a2b3c0d-2`,
			input: `a2b3c0d-2`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 2},
				{[]byte("b"), 3},
				{[]byte("c"), 0},
				{[]byte("d"), 1},
				{[]byte("-"), 2},
			}},
		},
		// escape tests
		{
			name:  `Simple escape`,
			input: `q\4\5`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("4"), 1},
				{[]byte("5"), 1},
			}},
		},
		{
			name:  `Escape repeat`,
			input: `q\45`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("4"), 5},
			}},
		},
		{
			name:  `Repeat tab (special character)`,
			input: "q\t3",
			expected: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("\t"), 3},
			}},
		},
		{
			name:  `Backslash repeat`,
			input: `q\\3`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("\\"), 3},
			}},
		},
		{
			name:  `Backslash escape`,
			input: `q\\\3`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("\\"), 1},
				{[]byte("3"), 1},
			}},
		},
		{
			name:  `Last backslash not escapes`,
			input: `qb2\`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("b"), 2},
				{[]byte("\\"), 1},
			}},
		},
		// tests with invalid repeat count position
		{
			name:     `Repeat times more then 9`,
			input:    `q12`,
			expected: GroupStorage{},
			err:      fmt.Errorf("invalid format q12, offset: 3, digit can't be the first symbol"),
		},
		{
			name:     `Starts from digit`,
			input:    `1a2`,
			expected: GroupStorage{},
			err:      fmt.Errorf("invalid format 1a2, offset: 1, digit can't be the first symbol"),
		},
		// tests with unicode
		{
			name:  `Single code-point emoji without skin tone modifier`,
			input: `a2👨2`,
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 2},
				{[]byte("👨"), 2},
			}},
		},
		{
			name:  `Multy code-point emoji (with skin tone modifier)`,
			input: `a3👨🏽2`,
			//       👨\xF0\x9F\x91\xA8 // person, code-point 1
			//                       🏽 \xF0\x9F\x8F\xBD // skin tone fitz-5, code-point 2
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 3},
				{[]byte("👨🏽"), 2},
			}},
		},
		{
			name:  `Multy code-point emoji (with skin tone modifier and zero width joiner)`,
			input: `b1👨🏾‍🚀2`,
			// 👨🏾‍🚀 =   "\xF0\x9F\x91\xA8\xF0\x9F\x8F\xBE\xE2\x80\x8B\xF0\x9F\x9A\x80",
			//       👨\xF0\x9F\x91\xA8 // person
			//                        🏽 \xF0\x9F\x8F\xBD // skin tone fitz-5
			//                                          \xE2\x80\x8B // zero width joiner
			//                                                    🚀\xF0\x9F\x9A\x80 //rocket
			expected: GroupStorage{[]repeatGroup{
				{[]byte("b"), 1},
				{[]byte("👨🏾‍🚀"), 2},
			}},
		},
		{
			name:  `Multy code-point emoji with several modifiers`,
			input: "a1e\u0301\u03012",
			// expected aé́é́
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte{0x65, 0xcc, 0x81, 0xcc, 0x81}, 2},
				//  e = 0x65
				//        ' = 0xcc, 0x81
				//                    ' = 0xcc, 0x81
			}},
		},
		{
			name:  `Multy code-point with letter modifier`,
			input: "a1e\u02EF2",
			// expected ae˯˯e˯˯
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte("e\u02EF"), 2},
			}},
		},
		{
			name:  `Non ascii digits in repeat count (MATHEMATICAL MONOSPACE DIGIT THREE)`,
			input: `a𝟹`,
			// input: "a\xF0\x9D\x9F\xB9",
			// it looks like a3 (@see https://www.fileformat.info/info/unicode/char/1d7f9/index.htm)
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte{0xF0, 0x9D, 0x9F, 0xB9}, 1},
			}},
		},
		{
			name:  `Repeat non ascii digit (MATHEMATICAL MONOSPACE DIGIT THREE)`,
			input: `a𝟹2`,
			// input: "a\xF0\x9D\x9F\xB92",
			// it looks like a3 (@see https://www.fileformat.info/info/unicode/char/1d7f9/index.htm)
			expected: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte{0xF0, 0x9D, 0x9F, 0xB9}, 2},
			}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseString(tt.input)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnpack(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		input    GroupStorage
		expected string
		err      error
	}{
		// trivial tests
		{
			name:     `Empty string`,
			input:    GroupStorage{},
			expected: "",
		},
		{
			name:     `Simple char`,
			input:    GroupStorage{[]repeatGroup{{[]byte("a"), 1}}},
			expected: `a`,
		},
		{
			name: `a2b3c0d-2`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 2},
				{[]byte("b"), 3},
				{[]byte("c"), 0},
				{[]byte("d"), 1},
				{[]byte("-"), 2},
			}},
			expected: `aabbbd--`,
		},
		{
			name: `Repeat tab (special character)`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("\t"), 3},
			}},
			expected: "q\t\t\t",
		},
		{
			name: `Backslash repeat`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("\\"), 3},
			}},
			expected: `q\\\`,
		},
		{
			name: `Last backslash not escapes`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("q"), 1},
				{[]byte("b"), 2},
				{[]byte("\\"), 1},
			}},
			expected: `qbb\`,
		},
		// tests with unicode
		{
			name: `Single code-point emoji without skin tone modifier`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 2},
				{[]byte("👨"), 2},
			}},
			expected: `aa👨👨`,
		},
		{
			name: `Multy code-point emoji (with skin tone modifier)`,
			//       👨\xF0\x9F\x91\xA8 // person, code-point 1
			//                       🏽 \xF0\x9F\x8F\xBD // skin tone fitz-5, code-point 2
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 3},
				{[]byte("👨🏽"), 2},
			}},
			expected: `aaa👨🏽👨🏽`,
		},
		{
			name: `Multy code-point emoji (with skin tone modifier and zero width joiner)`,
			// 👨🏾‍🚀 =   "\xF0\x9F\x91\xA8\xF0\x9F\x8F\xBE\xE2\x80\x8B\xF0\x9F\x9A\x80",
			//       👨\xF0\x9F\x91\xA8 // person
			//                        🏽 \xF0\x9F\x8F\xBD // skin tone fitz-5
			//                                          \xE2\x80\x8B // zero width joiner
			//                                                    🚀\xF0\x9F\x9A\x80 //rocket
			input: GroupStorage{[]repeatGroup{
				{[]byte("b"), 1},
				{[]byte("👨🏾‍🚀"), 2},
			}},
			expected: `b👨🏾‍🚀👨🏾‍🚀`,
		},
		{
			name: `Multy code-point emoji with several modifiers`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte{0x65, 0xcc, 0x81, 0xcc, 0x81}, 2},
				//             e = 0x65
				//                   ' = 0xcc, 0x81
				//                  			 ' = 0xcc, 0x81
			}},
			expected: `aé́é́`,
		},
		{
			name: `Multy code-point with letter modifier`,
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte("e\u02EF"), 2},
			}},
			expected: `ae˯e˯`,
		},
		{
			name: `Non ascii digits in repeat count (MATHEMATICAL MONOSPACE DIGIT THREE)`,
			// input: "a\xF0\x9D\x9F\xB9",
			// it looks like a3 (@see https://www.fileformat.info/info/unicode/char/1d7f9/index.htm)
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte{0xF0, 0x9D, 0x9F, 0xB9}, 1},
			}},
			expected: `a𝟹`,
		},
		{
			name: `Repeat non ascii digit (MATHEMATICAL MONOSPACE DIGIT THREE)`,
			// input: "a\xF0\x9D\x9F\xB92",
			// it looks like a3 (@see https://www.fileformat.info/info/unicode/char/1d7f9/index.htm)
			input: GroupStorage{[]repeatGroup{
				{[]byte("a"), 1},
				{[]byte{0xF0, 0x9D, 0x9F, 0xB9}, 2},
			}},
			expected: `a𝟹𝟹`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.Unpack()
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
