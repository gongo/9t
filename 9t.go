package main

import (
	"log"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	filenames := os.Args[1:]
	maxLength := maximumNameLength(filenames)

	for i, filename := range filenames {
		t, err := newTailer(filename, getColorCode(i), maxLength)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func() {
			t.do()
			wg.Done()
		}()
	}

	wg.Wait()
}
