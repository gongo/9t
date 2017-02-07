package ninetail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	ansiColorRegex = regexp.MustCompile("^\x1b\\[[0-9]+m(.*)\x1b\\[0m: (.*)$")
	defaultOutput  = colorableOutput
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

	tailers := NewTailers(names)

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

func TestTailerDoForLinesAreNiceAligned(t *testing.T) {
	// Stubbing output
	output := new(bytes.Buffer)
	colorableOutput = output
	defer revertDefault()

	dir, err := ioutil.TempDir("", "ninetail")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bases := []string{
		"test_tailer_do.log",
		"short.txt",
		"世界一かわいいよ.log",
	}
	files := make([]*os.File, len(bases))
	for i, name := range bases {
		// Create test file
		filename := filepath.Join(dir, name)
		f, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}
		files[i] = f
	}

	var wg sync.WaitGroup
	tailers := NewTailers(getFilenames(files))
	for _, t := range tailers {
		wg.Add(1)
		go func(t *Tailer) {
			t.Do()
			wg.Done()
		}(t)
	}

	// ^^;)
	interval := time.Tick(100 * time.Millisecond)
	<-interval
	fmt.Fprintf(files[0], "foobar\n")
	<-interval
	fmt.Fprintf(files[2], "その作業安全ですかと言える社風\n")
	<-interval
	fmt.Fprintf(files[1], "PPAP = Perfect PHP As PHP\n")
	<-interval

	for _, t := range tailers {
		t.Stop()
	}
	wg.Wait()

	expect := `  test_tailer_do.log: foobar
世界一かわいいよ.log: その作業安全ですかと言える社風
           short.txt: PPAP = Perfect PHP As PHP
`
	actual := stripAnsiColorCode(output.String())

	if expect != actual {
		t.Fatal("Incorrect align")
	}
}

func revertDefault() {
	colorableOutput = defaultOutput
}

func getFilenames(files []*os.File) []string {
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name()
	}
	return filenames
}

func stripAnsiColorCode(text string) string {
	stripped := ""

	for _, line := range strings.Split(text, "\n") {
		s := ansiColorRegex.FindStringSubmatch(line)

		if s != nil { // nil is empty line
			stripped += s[1] + ": " + s[2] + "\n"
		}
	}

	return stripped
}
