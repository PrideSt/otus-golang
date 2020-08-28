package pipe

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/top"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	goleak.VerifyTestMain(m)

	os.Exit(m.Run())
}

func TestNormalizeWord(t *testing.T) {
	normalizer := GetWordNormalizer()
	for _, tt := range [...]struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     `empty`,
			input:    ``,
			expected: []string(nil),
		},
		{
			name:     `one symbol`,
			input:    `a`,
			expected: []string{`a`},
		},
		{
			name:     `simple word`,
			input:    `word`,
			expected: []string{`word`},
		},
		{
			name:     `simple слово`,
			input:    `слово`,
			expected: []string{`слово`},
		},
		{
			name:     `dash`,
			input:    `-`,
			expected: []string(nil),
		},
		{
			name:     `dash first`,
			input:    `-one`,
			expected: []string{`one`},
		},
		{
			name:     `dash last`,
			input:    `one-`,
			expected: []string{`one`},
		},
		{
			name:     `dash in the middle`,
			input:    `какой-то`,
			expected: []string{`какой-то`},
		},
		{
			name:     `with punctuation last`,
			input:    `hello!`,
			expected: []string{`hello`},
		},
		{
			name:     `with punctuation in the middle`,
			input:    `hello,Masha`,
			expected: []string{`hello`, `masha`},
		},
		{
			name:     `case insensetive`,
			input:    `hElLo`,
			expected: []string{`hello`},
		},
		{
			name:     `with numbers`,
			input:    `i18n`,
			expected: []string{`i18n`},
		},
		{
			name:     `with special characters`,
			input:    "with\ttab",
			expected: []string{`with`, `tab`},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, normalizer(tt.input))
		})
	}
}

func TestGetWords(t *testing.T) {
	for _, tt := range [...]struct {
		name     string
		input    []top.FreqEntry
		expected []string
	}{
		{
			name:     `empty`,
			input:    []top.FreqEntry{},
			expected: []string{},
		},
		{
			name: `empty`,
			input: []top.FreqEntry{
				{"one", 1},
				{"two", 2},
				{"three", 3},
			},
			expected: []string{"one", "two", "three"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, GetWords(tt.input))
		})
	}
}

func TestTerminate(t *testing.T) {
	chTexts := make(chan string)
	chTerm := make(chan struct{})
	pushNewChunkInterval := 100 * time.Millisecond
	termInterval := 10 * pushNewChunkInterval

	wg := &sync.WaitGroup{}

	// ever push to text channel some text every 100ms
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(chTexts)

		word := `Фронт отклоняет циркулирующий поток. Конечно, нельзя не принять во внимание тот факт, что колебание семантически концентрирует абстрактный ритмический рисунок.`
		for {
			select {
			case chTexts <- word:
			case <-chTerm:
				return
			}
			time.Sleep(pushNewChunkInterval)
		}
	}()

	var result []top.FreqEntry
	wg.Add(1)
	go func() {
		defer wg.Done()
		// TopNInChan is blocking operation, it counts word entries using map synchronously
		// and we should wrap it in goroutine to can terminate after sleep in main thread
		result = TopNInChan(chTexts, 10, chTerm)
	}()

	time.Sleep(termInterval)
	close(chTerm)
	wg.Wait()

	// result can't be empty when already count some words in map and input channel closed
	// only when chTerm closed returns empty result
	require.Equal(t, []top.FreqEntry{}, result)
}
