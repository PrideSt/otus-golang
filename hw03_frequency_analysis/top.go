package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/topbuffer"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// chopString return first limit bytes string from s or full string if
// them length less then limit.
func chopString(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}

	return s
}

// getWordNormalizer returns function used for word normalization, returns proper words occurred in string.
func getWordNormalizer() func(s string) []string {
	pattern := `[a-zA-Zа-яА-Я0-9]+(?:-[a-zA-Zа-яА-Я0-9]+)*`
	re := regexp.MustCompile(pattern)

	return func(s string) []string {
		lowerStr := strings.ToLower(s)

		return re.FindAllString(lowerStr, -1)
	}
}

// operate runs goroutine which apply f function to every entry in input stream and write result to output stream.
func operate(f func(s string) []string, in <-chan string, chTerminate <-chan struct{}, wg *sync.WaitGroup) <-chan string {
	out := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		for {
			select {
			case text, ok := <-in:
				if !ok {
					logger.Println("all done, terminate operation")

					return
				}
				logger.Printf("input \"%s...\"\n", chopString(text, 10))

				for i, chunk := range f(text) {
					select {
					case out <- chunk:
						logger.Printf("out[%d]: \"%s...\"\n", i, chopString(chunk, 10))
					case <-chTerminate:
						logger.Println("terminate operation gracefully")

						return
					}
				}
			case <-chTerminate:
				logger.Println("terminate operation gracefully")

				return
			}
		}
	}()

	return out
}

// findTopN pass every antry from ferq dictionary (map of word -> count in text) through topbuffer.Interface- ordered buffern with
// topResultsCnt elements. Returns ordered slice of pair {word, freq}.
func findTopN(dict map[string]int, topResultsCnt int, chTerminate <-chan struct{}) []topbuffer.FreqEntry {
	top := topbuffer.New(topResultsCnt, func(lhs, rhs topbuffer.FreqEntry) bool {
		return lhs.Freq > rhs.Freq
	})

	logger.Println("Find top, fill top buffer")
	i := 0
	for w, fr := range dict {
		logger.Println(w, "->", fr)
		if i%100 == 0 {
			// we should terminate while finding top, check it every 100 iterations
			select {
			case <-chTerminate:
				logger.Println("terminate find-top-n")

				return []topbuffer.FreqEntry{}
			default:
			}
		}
		top.Add(topbuffer.FreqEntry{Word: w, Freq: fr})
		i++
	}

	return top.Get()
}

// getTopWords converts list of {word, freq} pairs to word slice.
func getTopWords(r []topbuffer.FreqEntry) []string {
	results := make([]string, len(r))
	for i, pair := range r {
		results[i] = pair.Word
	}

	return results
}

// countWords sum how much times words in chNormWords occures.
func countWords(chNormWords <-chan string, chTerminate <-chan struct{}) map[string]int {
	dict := make(map[string]int)
	for {
		select {
		case word, ok := <-chNormWords:
			if !ok {
				logger.Println("All words counted, terminate freq-counter")

				return dict
			}
			logger.Printf("increase count of word: %q\n", word)
			dict[word]++
		case <-chTerminate:
			logger.Println("gracefully freq-counter termination")

			return nil
		}
	}
}

// TopN return topLen the most frequencies words in input.
func TopN(input string, topLen int, chTerminate <-chan struct{}) []topbuffer.FreqEntry {
	chText := make(chan string, 32)

	wg := sync.WaitGroup{}
	// if chTerminate closed we must wait running goroutins
	defer wg.Wait()

	// split input text on chuncks, maxLen must be increased in real cases (text can be long)
	chTextChunks := operate(
		func(delim string, maxLen int) func(s string) []string {
			return func(input string) []string {
				chunkCnt := len(input)/maxLen + 1
				return strings.SplitN(input, delim, chunkCnt)
			}
		}(" ", 16),
		chText,
		chTerminate,
		&wg,
	)

	// split text chunks on words
	chDryWords := operate(
		strings.Fields,
		chTextChunks,
		chTerminate,
		&wg,
	)

	// convert word in normal form
	chNormWords := operate(
		getWordNormalizer(),
		chDryWords,
		chTerminate,
		&wg,
	)

	chText <- input
	close(chText)

	// count words frequencies in main thread, map requires serialization
	dict := countWords(chNormWords, chTerminate)

	// pass map through TopBuffer and find topLen the most frequencies words
	top := findTopN(dict, topLen, chTerminate)

	return top
}

// Top10 returns 10 the most frequencies words in s.
func Top10(s string, chTerminate <-chan struct{}) []string {
	return getTopWords(TopN(s, 10, chTerminate))
}
