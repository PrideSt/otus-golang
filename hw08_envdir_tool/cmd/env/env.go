package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	env := os.Environ()
	sort.Sort(sort.StringSlice(env))
	for _, pair := range env {
		fmt.Println(pair)
	}
}
