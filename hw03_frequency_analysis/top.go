package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/topbuffer"
)

const isVerbose bool = false

// chopString return first limit bytes string from s or full string if
// them length less then limit.
func chopString(s string, limit int) string {
	if len(s) > limit {
		return s[:limit]
	}

	return s
}

// runTextChunker starts goroutine which divide input string by chunks with size less then maxLen and
// and write chunks to out channel.
func runTextChunker(delim string, maxLen int, in <-chan string, out chan<- string, chTerminate <-chan struct{}, wg *sync.WaitGroup, log *log.Logger) {
	defer wg.Done()
	defer close(out)

	for {
		select {
		case text, ok := <-in:
			if !ok {
				log.Println("All texts splitted, terminate text-chunker")

				return
			}
			log.Printf("let's chunck the text: \"%s...\"\n", chopString(text, 10))

			chunkCnt := len(text)/maxLen + 1
			chunks := strings.SplitN(text, delim, chunkCnt)

			for i, chunk := range chunks {
				select {
				case out <- chunk:
					log.Printf("write the chunk #%d: \"%s...\"\n", i, chopString(chunk, 10))
				case <-chTerminate:
					log.Println("gracefully text-chunker termination")

					return
				}
			}
		case <-chTerminate:
			log.Println("gracefully text-chunker termination")

			return
		}
	}
}

// runWordSplitter starts goroutine which extracts words from string passed to in channel
// and write word as is to out channel.
func runWordSplitter(in <-chan string, out chan<- string, chTerminate <-chan struct{}, wg *sync.WaitGroup, log *log.Logger) {
	defer wg.Done()
	defer close(out)

	for {
		select {
		case text, ok := <-in:
			if !ok {
				log.Println("All words splitted, terminate word-splitter")

				return
			}
			log.Printf("input text appears: \"%s...\", split it on words\n", chopString(text, 10))
			for _, word := range strings.Fields(text) {
				select {
				case out <- word:
					log.Printf("writed raw word: %q\n", word)
				case <-chTerminate:
					log.Println("gracefully word-splitter termination")

					return
				}
			}
		case <-chTerminate:
			log.Println("gracefully word-splitter termination")

			return
		}
	}
}

// runWordNormalizer starts goroutine to converts raw words to notmal form using function norm, apply them
// to every word from in channel and pass normalization result to out channel.
func runWordNormalizer(norm func(s string) []string, in <-chan string, out chan<- string, chTerminate <-chan struct{}, wg *sync.WaitGroup, log *log.Logger) {
	defer wg.Done()
	defer close(out)

	for {
		select {
		case word, ok := <-in:
			if !ok {
				defer log.Println("All words normalized, terminate word-normalizer")

				return
			}
			for i, normWord := range norm(word) {
				log.Printf("normolize word %q part %d -> %q\n", word, i, normWord)

				select {
				case out <- normWord:
					log.Printf("writed normal word : %q\n", normWord)
				case <-chTerminate:
					log.Println("gracefully word-normalizer termination")

					return
				}
			}
		case <-chTerminate:
			log.Println("gracefully word-normalizer termination")

			return
		}
	}
}

// runFreqCounter starts goroutine, get normalized word into in channel and increase frequencies of
// this word into internal map dict. When in channel closed and all words treated close done
// channel and unlock main thred for search N max (the most frequencies) elements.
func runFreqCounter(in <-chan string, out chan<- map[string]int, chTerminate <-chan struct{}, wg *sync.WaitGroup, log *log.Logger) {
	defer wg.Done()
	defer close(out)

	dict := make(map[string]int)
	// we use index map in single thread (single-writer) and shouldn't use locks
	for {
		select {
		case word, ok := <-in:
			if !ok {
				log.Println("All words counted, terminate freq-counter")
				out <- dict

				return
			}
			log.Printf("increase count of word: %q\n", word)
			dict[word]++
		case <-chTerminate:
			log.Println("gracefully freq-counter termination")

			return
		}
	}
}

// runSygnalListener subscribes on chSygnal channel and close chTerminate channel on any income sygnal.
func runSygnalListener(wg *sync.WaitGroup, chTerminate chan struct{}, chSygnal <-chan os.Signal, log *log.Logger) {
	defer wg.Done()
	// close terminate channel can only this listener, but it can't write to them
	// main thread can write, but not close
	defer close(chTerminate)

	_, ok := <-chSygnal
	if ok {
		log.Println("terminate sygnal received")
	} else {
		log.Println("terminate sygnal-listener")
	}
}

// findTopN pass every antry from ferq dictionary (map of word -> count in text) through topbuffer.Interface- ordered buffern with
// topResultsCnt elements. Returns ordered slice of pair {word, freq}.
func findTopN(dict map[string]int, topResultsCnt int, chTerminate <-chan struct{}, log *log.Logger) ([]topbuffer.FreqEntry, error) {
	top := topbuffer.New(topResultsCnt, func(lhs, rhs topbuffer.FreqEntry) bool {
		return lhs.Freq > rhs.Freq
	})

	log.Println("Find top, fill top buffer")
	i := 0
	for w, fr := range dict {
		log.Println(w, "->", fr)
		if i%100 == 0 {
			// we should terminate while finding top
			select {
			case <-chTerminate:
				log.Println("terminate find-top-n")

				return []topbuffer.FreqEntry{}, fmt.Errorf("findTopN terminated")
			default:
			}
		}
		top.Add(topbuffer.FreqEntry{Word: w, Freq: fr})
		i++
	}

	return top.Get(), nil
}

// getTopWords converts list of {word, freq} pairs to word slice.
func getTopWords(r []topbuffer.FreqEntry) []string {
	results := make([]string, len(r))
	for i, pair := range r {
		results[i] = pair.Word
	}

	return results
}

// getWordNormalizer returns normalize function, allow reuse compiled regexp.
func getWordNormalizer() func(s string) []string {
	pattern := `[a-zA-Zа-яА-Я0-9]+(?:-[a-zA-Zа-яА-Я0-9]+)*`
	re := regexp.MustCompile(pattern)

	return func(s string) []string {
		lowerStr := strings.ToLower(s)

		return re.FindAllString(lowerStr, -1)
	}
}

// TopN return topLen the most frequencies words in input.
func TopN(input string, topLen int) []topbuffer.FreqEntry {
	var plogger *log.Logger

	if isVerbose {
		plogger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		plogger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	plogger.Println("my pid is:", os.Getpid())

	chText := make(chan string, 32)
	chTextChunks := make(chan string, 32)
	chDryWords := make(chan string, 32)
	chNormWords := make(chan string, 32)
	chDict := make(chan map[string]int, 1)

	chSygnal := make(chan os.Signal, 1)
	signal.Notify(chSygnal, syscall.SIGTERM, syscall.SIGINT)

	chTerminate := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go runSygnalListener(&wg, chTerminate, chSygnal, plogger)

	wg.Add(4)
	go runTextChunker(" ", 16, chText, chTextChunks, chTerminate, &wg, plogger)
	go runWordSplitter(chTextChunks, chDryWords, chTerminate, &wg, plogger)
	go runWordNormalizer(getWordNormalizer(), chDryWords, chNormWords, chTerminate, &wg, plogger)
	go runFreqCounter(chNormWords, chDict, chTerminate, &wg, plogger)

	chText <- input
	close(chText)

	plogger.Println("Wait all wards counted...")
	dict := <-chDict

	top, err := findTopN(dict, topLen, chTerminate, plogger)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// close sygnal channel and terminate sygnal listener goroutine
	close(chSygnal)

	// when we catch a sygnal wait till all goroutins closed
	wg.Wait()

	return top
}

// Top10 returns 10 the most frequencies words in s.
func Top10(s string) []string {
	return getTopWords(TopN(s, 10))
}
