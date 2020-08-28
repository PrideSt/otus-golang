package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/pipe"
	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/top"
)

// Top10 returns 10 the most frequencies words in s.
func Top10(s string, chTerminate <-chan struct{}) []string {
	return pipe.GetWords(TopN(s, 10, chTerminate))
}

// TopN return topLen the most frequencies words in input string.
func TopN(s string, topLen int, chTerminate <-chan struct{}) []top.FreqEntry {
	normalizer := pipe.GetWordNormalizer()
	wordCnt := make(map[string]int)
	for _, str := range strings.Fields(s) {
		for _, ss := range normalizer(str) {
			wordCnt[ss]++
		}
	}

	topWords := make(top.FreqSlice, 0, len(wordCnt))
	for word, cnt := range wordCnt {
		topWords = append(topWords, top.FreqEntry{
			Word: word,
			Freq: cnt,
		})
	}

	sort.Sort(sort.Reverse(topWords))

	if len(topWords) < topLen {
		return topWords
	}

	return topWords[:topLen]
}
