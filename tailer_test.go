package ninetail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ActiveState/tail"
)

var (
	defaultSeekInfo = seekInfoOnStart
	defaultOutput   = colorableOutput
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
	// Change seek to start
	seekInfoOnStart = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_SET}
	defer revertDefault()

	// Create test file
	testfile, err := ioutil.TempFile(os.TempDir(), "TestTailerDo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testfile.Name())

	testfile.WriteString("foobar")

	tailer, err := newTailer(testfile.Name(), 31, 10)
	if err != nil {
		t.Fatal(err)
	}

	go tailer.Do()

	timeout := time.After(300 * time.Millisecond)
	<-timeout
	tailer.Stop()

	expect := fmt.Sprintf(
		"\x1b[31m%s\x1b[0m: foobar\n",
		path.Base(testfile.Name()),
	)

	if expect != output.String() {
		t.Fatal("hogehoge")
	}
}

func revertDefault() {
	colorableOutput = defaultOutput
	seekInfoOnStart = defaultSeekInfo
}
