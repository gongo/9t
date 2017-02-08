package ninetail

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/mattn/go-colorable"
)

func TestNewTailers(t *testing.T) {
	names := []string{
		"/very/very/long/path/but/base/is/short",
		"/path/to/message",
		"/p/t/very_very_long_base",
		"/p/var/log/production.log",
		"/p/var/log/development.log",
		"/p/var/log/test.log",
	}

	tailers, err := NewTailers(names)
	if err != nil {
		t.Fatal(err)
	}

	if len(tailers) != len(names) {
		t.Fatalf("Incorrect: tailers count expect(%d) actual(%d)", len(names), len(tailers))
	}

	if tailers[0].colorCode != ansiColorCodes[0] {
		t.Fatal("Incorrect color code at 1st file")
	}

	if tailers[1].colorCode != ansiColorCodes[1] {
		t.Fatal("Incorrect color code at 2nd file")
	}

	if tailers[5].colorCode != ansiColorCodes[0] { // Return to first color code
		t.Fatal("Incorrect color code at 6th file")
	}
}

func TestTailerDo(t *testing.T) {
	file, err := ioutil.TempFile("", "ninetail_tailer_do")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	basename := filepath.Base(file.Name())
	maxLength := len(basename) + 20 // 20 = padding

	tailer, err := newTailer(file.Name(), 32, maxLength) // 32 = green
	if err != nil {
		t.Fatal(err)
	}

	output := new(bytes.Buffer)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(t *Tailer, output io.Writer) {
		tailer.Do(colorable.NewNonColorable(output))
		wg.Done()
	}(tailer, output)

	// Simulate to `echo line >> ninetail_tailer_do`
	interval := time.Tick(100 * time.Millisecond)
	<-interval
	fmt.Fprint(file, "line\n")
	<-interval

	tailer.Stop()
	wg.Wait()

	expect := fmt.Sprintf("                    %s: line\n", basename)
	actual := output.String()

	if expect != actual {
		t.Fatal("Bad padding or line text")
	}
}
