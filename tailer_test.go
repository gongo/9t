package ninetail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	defaultOutput = colorableOutput
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

	if tailers[0].maxWidth != 19 { // len("very_very_long_base")
		t.Fatalf("Incorrect: maximum name length: expect(%d) actual(%d)", 19, tailers[0].maxWidth)
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
	// Stubbing output
	output := new(bytes.Buffer)
	colorableOutput = output
	defer revertDefault()

	// Create test file
	testfile, err := ioutil.TempFile(os.TempDir(), "TestTailerDo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testfile.Name())

	tailer, err := newTailer(testfile.Name(), 31, 10)
	if err != nil {
		t.Fatal(err)
	}

	go tailer.Do()

	interval := time.Tick(100 * time.Millisecond)
	<-interval
	testfile.WriteString("foobar\n")
	<-interval
	tailer.Stop()

	expect := fmt.Sprintf(
		"\x1b[%dm%*s\x1b[0m: foobar\n",
		31, // eq 2nd args on newTailer()
		10, // eq 3rd args on newTailer()
		filepath.Base(testfile.Name()),
	)

	if expect != output.String() {
		t.Fatal("Incorrect display")
	}
}

func revertDefault() {
	colorableOutput = defaultOutput
}

func TestMaximumNameLength(t *testing.T) {
	ns := []struct {
		name   string
		length int
	}{
		{"a", 1},
		{"ab", 2},
		{"世界一かわいいよ", 8},
		{"/p/t/very_very_long_base", 19}}
	names := make([]string, 4)
	for _, n := range ns {
		names = append(names, n.name)
		cl := maximumNameLength(names)
		if n.length != cl {
			t.Fatalf("Incorrect: Maximum name length: expect(%d) actual(%d)", n.length, cl)
		}
	}
}
