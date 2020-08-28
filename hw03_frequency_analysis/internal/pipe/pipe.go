package pipe

import (
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/top"
)

// GetWords converts list of {word, freq} pairs to word slice.
func GetWords(r []top.FreqEntry) []string {
	results := make([]string, len(r))
	for i, pair := range r {
		results[i] = pair.Word
	}

	return results
}

// chopString return first limit bytes string from s or full string if
// them length less then limit.
func chopStringN(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}

	return s
}

// chopString10 is shortcut for chopStringN(s, 10).
func chopString10(s string) string {
	return chopStringN(s, 10)
}

// GetWordNormalizer returns function used for word normalization, returns proper words occurred in string.
func GetWordNormalizer() func(s string) []string {
	pattern := `[a-zA-Zа-яА-Я0-9]+(?:-[a-zA-Zа-яА-Я0-9]+)*`
	re := regexp.MustCompile(pattern)

	return func(s string) []string {
		lowerStr := strings.ToLower(s)
		time.Sleep(time.Second)

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
					log.Println("all done, terminate operation")

					return
				}
				log.Printf("input \"%s...\"\n", chopString10(text))

				for i, chunk := range f(text) {
					select {
					case out <- chunk:
						log.Printf("out[%d]: \"%s...\"\n", i, chopString10(chunk))
					case <-chTerminate:
						log.Println("terminate operation gracefully")

						return
					}
				}
			case <-chTerminate:
				log.Println("terminate operation gracefully")

				return
			}
		}
	}()

	return out
}

// countWords sum how much times words in chNormWords occures.
func countWords(chNormWords <-chan string, chTerminate <-chan struct{}) map[string]int {
	dict := make(map[string]int)
	for {
		select {
		case word, ok := <-chNormWords:
			if !ok {
				log.Println("All words counted, terminate freq-counter")

				return dict
			}
			log.Printf("increase count of word: %q\n", word)
			dict[word]++
		case <-chTerminate:
			log.Println("gracefully freq-counter termination")

			return nil
		}
	}
}

// TopNInChan return topLen the most frequencies words in chTextChunks channel.
func TopNInChan(chTextChunks <-chan string, topLen int, chTerminate <-chan struct{}) []top.FreqEntry {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// split text chunks on words
	chDryWords := operate(
		strings.Fields,
		chTextChunks,
		chTerminate,
		wg,
	)

	// convert word in normal form
	chNormWords := operate(
		GetWordNormalizer(),
		chDryWords,
		chTerminate,
		wg,
	)

	// count words frequencies in main thread, map requires serialization
	dict := countWords(chNormWords, chTerminate)

	// pass map through Buffer and find topLen the most frequencies words
	return top.FindTopN(dict, topLen, chTerminate)
}

// TopN return topLen the most frequencies words in input string.
func TopN(input string, topLen int, chTerminate <-chan struct{}) []top.FreqEntry {
	delim := " "
	maxLen := 4096
	chunkCnt := len(input)/maxLen + 1
	chTextChunks := make(chan string)

	go func() {
		defer close(chTextChunks)
		for _, chunk := range strings.SplitN(input, delim, chunkCnt) {
			select {
			case <-chTerminate:
				log.Printf("terminate input channel\n")
				return
			default:
			}
			select {
			case <-chTerminate:
				log.Printf("terminate input channel\n")
				return
			case chTextChunks <- chunk:
				log.Printf("write text chunk %q to input stream\n", chopString10(chunk))
			}
		}
	}()

	result := TopNInChan(chTextChunks, topLen, chTerminate)

	select {
	case <-chTextChunks:
	case <-chTerminate:
	}

	return result
}

// Top10 returns 10 the most frequencies words in s.
func Top10(s string, chTerminate <-chan struct{}) []string {
	return GetWords(TopN(s, 10, chTerminate))
}
