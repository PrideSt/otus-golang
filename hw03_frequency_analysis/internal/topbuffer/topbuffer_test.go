package topbuffer

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

var cmpLt = func(lhs, rhs FreqEntry) bool {
	return lhs.Freq < rhs.Freq
}

var cmpGt = func(lhs, rhs FreqEntry) bool {
	return lhs.Freq != rhs.Freq && !cmpLt(lhs, rhs)
}

type insertValAndPos struct {
	val FreqEntry
	pos int
}

// sort.Interface need for brnchmark
type freqSlice []FreqEntry

func (p freqSlice) Len() int           { return len(p) }
func (p freqSlice) Less(i, j int) bool { return p[i].Freq < p[j].Freq }
func (p freqSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func TestAdd(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		len      int
		less     func(lhs, rhs FreqEntry) bool
		input    []insertValAndPos
		expected []FreqEntry
	}{
		{
			name:     "empty",
			len:      3,
			less:     cmpLt,
			input:    []insertValAndPos{},
			expected: []FreqEntry{},
		},
		{
			name: "simple lt",
			len:  3,
			less: cmpLt,
			input: []insertValAndPos{
				{FreqEntry{"two", 2}, 0},
				{FreqEntry{"one", 1}, 0},
				{FreqEntry{"ten", 10}, 2},
			},
			expected: []FreqEntry{
				{"one", 1},
				{"two", 2},
				{"ten", 10},
			},
		},
		{
			name: "simple gt",
			len:  3,
			less: cmpGt,
			input: []insertValAndPos{
				{FreqEntry{"two", 2}, 0},
				{FreqEntry{"one", 1}, 1},
				{FreqEntry{"ten", 10}, 0},
			},
			expected: []FreqEntry{
				{"ten", 10},
				{"two", 2},
				{"one", 1},
			},
		},
		{
			name: "simple overflow",
			len:  3,
			less: cmpGt,
			input: []insertValAndPos{
				{FreqEntry{"five", 5}, 0},
				{FreqEntry{"three", 3}, 1},
				{FreqEntry{"one", 1}, 2},
				{FreqEntry{"zero", 0}, -1},
			},
			expected: []FreqEntry{
				{"five", 5},
				{"three", 3},
				{"one", 1},
			},
		},
		{
			name: "replace last with min",
			len:  3,
			less: cmpGt,
			input: []insertValAndPos{
				{FreqEntry{"five", 5}, 0},
				{FreqEntry{"three", 3}, 1},
				{FreqEntry{"one", 1}, 2},
				{FreqEntry{"two", 2}, 2},
			},
			expected: []FreqEntry{
				{"five", 5},
				{"three", 3},
				{"two", 2},
			},
		},
		{
			name: "replace last with max",
			len:  3,
			less: cmpGt,
			input: []insertValAndPos{
				{FreqEntry{"five", 5}, 0},
				{FreqEntry{"three", 3}, 1},
				{FreqEntry{"one", 1}, 2},
				{FreqEntry{"seven", 7}, 0},
			},
			expected: []FreqEntry{
				{"seven", 7},
				{"five", 5},
				{"three", 3},
			},
		},
		{
			name: "not full",
			len:  3,
			less: cmpGt,
			input: []insertValAndPos{
				{FreqEntry{"five", 5}, 0},
				{FreqEntry{"three", 3}, 1},
			},
			expected: []FreqEntry{
				{"five", 5},
				{"three", 3},
			},
		},
		{
			name: "not full reorder",
			len:  3,
			less: cmpGt,
			input: []insertValAndPos{
				{FreqEntry{"three", 3}, 0},
				{FreqEntry{"five", 5}, 0},
			},
			expected: []FreqEntry{
				{"five", 5},
				{"three", 3},
			},
		},
		{
			name: "custom comparator less",
			len:  3,
			less: cmpLt,
			input: []insertValAndPos{
				{FreqEntry{"five", 5}, 0},
				{FreqEntry{"three", 3}, 0},
				{FreqEntry{"one", 1}, 0},
				{FreqEntry{"seven", 7}, -1},
			},
			expected: []FreqEntry{
				{"one", 1},
				{"three", 3},
				{"five", 5},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			top := New(tt.len, tt.less)

			for _, pair := range tt.input {
				result := top.Add(pair.val)
				require.Equalf(t, pair.pos, result, "add value %d has wrong position %d, expected %d", pair.val, result, pair.pos)
			}

			require.Equal(t, tt.expected, top.Get())
		})
	}
}

func generateFreqMap(sz int, maxFreq int) map[string]int {
	m := make(map[string]int, sz)

	for i := 0; i < sz; i++ {
		m[fmt.Sprintf("word_%d", i)] = rand.Intn(maxFreq)
	}

	return m
}

func benchTopNFromM(t *testing.B, resultCnt int, dictSz int) {
	// use same freq dictionary for both tests
	m := generateFreqMap(dictSz, resultCnt*10)

	// not measure time on creating map
	t.ResetTimer()

	t.Run("measure slice", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			slice := make(freqSlice, len(m))

			// copy map to slice
			for word, freq := range m {
				slice = append(slice, FreqEntry{word, freq})
			}

			sort.Sort(slice)
			// result := slice[:resultCnt]
		}
	})
	t.Run("measure TopBuffer", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			tb := New(resultCnt, cmpGt)

			for word, freq := range m {
				tb.Add(FreqEntry{word, freq})
			}
			// result := tb.Get()
		}
	})
}

// BenchmarkTop10Of100/measure_slice-8         	   90361	     13306 ns/op	    7584 B/op	       3 allocs/op
// BenchmarkTop10Of100/measure_TopBuffer-8     	  236230	      5165 ns/op	     240 B/op	       1 allocs/op
func BenchmarkTop10Of100(t *testing.B) {
	benchTopNFromM(t, 10, 100)
}

// BenchmarkTop10Of1K/measure_slice-8          	    8028	    149028 ns/op	   73760 B/op	       3 allocs/op
// BenchmarkTop10Of1K/measure_TopBuffer-8      	   39016	     30559 ns/op	     240 B/op	       1 allocs/op
func BenchmarkTop10Of1K(t *testing.B) {
	benchTopNFromM(t, 10, 1000)
}

// BenchmarkTop10Of10K/measure_slice-8         	     679	   1757990 ns/op	 1417253 B/op	       5 allocs/op
// BenchmarkTop10Of10K/measure_TopBuffer-8     	    4742	    251064 ns/op	     240 B/op	       1 allocs/op
func BenchmarkTop10Of10K(t *testing.B) {
	benchTopNFromM(t, 10, 10000)
}

// BenchmarkTop100Of10K/measure_slice-8        	     525	   2271061 ns/op	 1417249 B/op	       5 allocs/op
// BenchmarkTop100Of10K/measure_TopBuffer-8    	    1930	    604759 ns/op	    2688 B/op	       1 allocs/op
func BenchmarkTop100Of10K(t *testing.B) {
	benchTopNFromM(t, 100, 10000)
}
