package main

import (
	"os"
	"sync"

	"github.com/gongo/9t"
)

func main() {
	var wg sync.WaitGroup
	filenames := os.Args[1:]

	for _, t := range ninetail.NewTailers(filenames) {
		wg.Add(1)
		go func(t *ninetail.Tailer) {
			t.Do()
			wg.Done()
		}(t)
	}

	wg.Wait()
}
