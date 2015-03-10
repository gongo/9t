package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/ActiveState/tail"
)

func maximumNameLength(filenames []string) int {
	length := 0
	for _, name := range filenames {
		if len(name) > length {
			length = len(name)
		}
	}
	return length
}

func tailer(t *tail.Tail, maxLength int, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range t.Lines {
		fmt.Printf("%*s: %s\n", maxLength, path.Base(t.Filename), line.Text)
	}
}

func main() {
	var wg sync.WaitGroup
	offset := &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	maxLength := maximumNameLength(os.Args[1:])

	for _, filename := range os.Args[1:] {
		t, err := tail.TailFile(filename, tail.Config{
			Follow:   true,
			Location: offset,
			Logger:   tail.DiscardingLogger,
		})

		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go tailer(t, maxLength, &wg)
	}

	wg.Wait()
}
