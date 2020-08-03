package main
// package hw03_frequency_analysis //nolint:golint,stylecheck

import(
	"strings"
	"log"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/txt"
)

func runTextChunker(delim string, maxLen int, in <-chan string, out chan<- string, chDone <-chan struct{}){
	defer close(out)
	
	select {
	case text := <- in:
		log.Println("let's chunck the text:", text[:10], "..., split on words")
		chunker := txt.Chunker(text, delim, maxLen)
		for i, chunk := range txt.Chunks(chunker) {
			select {
			case out <- chunk:
				log.Printf("write the chunk #%d: %s...\n", i, chunk[:10])
			case <-chDone:
				log.Println("gracefully termination")
				return
			}
		}
	case <- chDone:
		log.Println("gracefully termination")
		return
	}
}

func runWordSplitter(delim string, in <-chan string, out chan<- string, chDone <-chan struct{}){
	defer close(out)

	select {
	case text := <- in:
		log.Println("input text appears: ", text[:10], "..., split on words")
		for _, word := range strings.Split(text, delim) {
			select {
			case out <- word:
				log.Print("writed raw word:", word)
			case <-chDone:
				log.Println("gracefully termination")
				return
			}
		}
	case <- chDone:
		log.Println("gracefully termination")
		return
	}
}

// func Top10(_ string) []string {
func main() {
	input := "one two two three nine six one one six six six six"
	
	chText := make(chan string, 32)

	chDryWords := make(chan string, 32)
	// chNormWords := make(chan string)
	chDone := make(chan struct{})

	go runWordSplitter(chText, chDryWords, chDone)

	
	chText <- input
	close(chText)

	index := make(map[string]int)
	for word := range chDryWords {
		log.Println("catch word:", word)
		index[word]++
	}

	for k, v := range index {
		log.Println(k, "->", v)
	}
	// выгребаем слова из текста
	// нормализируем слова
	// считаем слова
	// ? можно ли использовать кучу
	// ели кучу использовать нельзя, то нужно найти топ 10 из мапки это partial_sort
	// отдать топ 10
	// Place your code here

	// return nil
}
