package topbuffer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type valPos struct {
	val int
	pos int
}

var cmpGt = func(lhs, rhs int) bool {
	return lhs > rhs
}

var cmpLt = func(lhs, rhs int) bool {
	return rhs > lhs
}

func TestAdd(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		size     int
		fgt      func(lhs, rhs int) bool
		input    []valPos
		expected []int
	}{
		{
			name: "simple",
			size: 3,
			fgt:  cmpGt,
			input: []valPos{
				{5, 0},
				{3, 1},
				{1, 2},
			},
			expected: []int{5, 3, 1},
		},
		{
			name: "simple overflow",
			size: 3,
			fgt:  cmpGt,
			input: []valPos{
				{5, 0},
				{3, 1},
				{1, 2},
				{0, -1},
			},
			expected: []int{5, 3, 1},
		},
		{
			name: "replace last with min",
			size: 3,
			fgt:  cmpGt,
			input: []valPos{
				{5, 0},
				{3, 1},
				{1, 2},
				{2, 2},
			},
			expected: []int{5, 3, 2},
		},
		{
			name: "replace last with max",
			size: 3,
			fgt:  cmpGt,
			input: []valPos{
				{5, 0},
				{3, 1},
				{1, 2},
				{7, 0},
			},
			expected: []int{7, 5, 3},
		},
		{
			name: "not full",
			size: 3,
			fgt:  cmpGt,
			input: []valPos{
				{5, 0},
				{3, 1},
			},
			expected: []int{5, 3},
		},
		{
			name: "not full reorder",
			size: 3,
			fgt:  cmpGt,
			input: []valPos{
				{3, 0},
				{5, 0},
			},
			expected: []int{5, 3},
		},
		{
			name: "custom comparator less",
			size: 3,
			fgt:  cmpLt,
			input: []valPos{
				{5, 0},
				{3, 0},
				{1, 0},
				{7, -1},
			},
			expected: []int{1, 3, 5},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			top := TopBuffer{
				buffer: []int{},
				size:   tt.size,
				fgt:    tt.fgt,
			}

			for _, pair := range tt.input {
				result := top.Add(pair.val)
				require.Equalf(t, pair.pos, result, "add value %d has wrong position %d, expected %d", pair.val, result, pair.pos)
			}

			require.Equal(t, tt.expected, top.Get())
		})
	}
}
