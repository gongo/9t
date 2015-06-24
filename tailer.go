package ninetail

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ActiveState/tail"
	"github.com/mattn/go-colorable"
)

var (
	// red, green, yellow, magenta, cyan
	ansiColorCodes  = [...]int{31, 32, 33, 35, 36}
	seekInfoOnStart = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	colorableOutput = colorable.NewColorableStdout()
)

type Tailer struct {
	*tail.Tail
	colorCode int
	maxWidth  int
}

func NewTailers(filenames []string) []*Tailer {
	maxLength := maximumNameLength(filenames)
	ts := make([]*Tailer, len(filenames))

	for i, filename := range filenames {
		t, err := newTailer(filename, getColorCode(i), maxLength)
		if err != nil {
			log.Fatal(err)
		}

		ts[i] = t
	}

	return ts
}

func newTailer(filename string, colorCode int, maxWidth int) (*Tailer, error) {
	t, err := tail.TailFile(filename, tail.Config{
		Follow:   true,
		Location: seekInfoOnStart,
		Logger:   tail.DiscardingLogger,
	})

	if err != nil {
		return nil, err
	}

	return &Tailer{
		Tail:      t,
		colorCode: colorCode,
		maxWidth:  maxWidth,
	}, nil
}

func (t Tailer) Do() {
	for line := range t.Lines {
		fmt.Fprintf(colorableOutput, "\x1b[%dm%*s\x1b[0m: %s\n", t.colorCode, t.maxWidth, t.name(), line.Text)
	}
}

func (t Tailer) name() string {
	return filepath.Base(t.Filename)
}

func getColorCode(index int) int {
	return ansiColorCodes[index%len(ansiColorCodes)]
}

func maximumNameLength(filenames []string) int {
	max := 0
	for _, name := range filenames {
		base := filepath.Base(name)
		if len(base) > max {
			max = len(base)
		}
	}
	return max
}
