package ninetail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/mattn/go-colorable"
)

func TestRun(t *testing.T) {
	tfs, err := newTestFileSet([]string{
		"php.txt",
		"吾輩は猫である.txt",
		"i_am_a_CAT.txt",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer tfs.removeAll()

	tailers, err := NewTailers(tfs.getFilenames())
	if err != nil {
		t.Fatal(err)
	}

	output := new(bytes.Buffer)
	target := &NineTail{
		output:  colorable.NewNonColorable(output),
		tailers: tailers,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(n *NineTail) {
		n.Run()
		wg.Done()
	}(target)

	interval := time.Tick(100 * time.Millisecond)
	<-interval
	tfs.writeString(0, "PHP(PHP: Hypertext Preprocessor)\n")
	<-interval
	tfs.writeString(2, "I am a cat. As yet I have no name.\n")
	<-interval
	tfs.writeString(1, "吾輩は猫である。名前はまだ無い。\n")
	<-interval

	for _, t := range target.tailers {
		t.Stop()
	}
	wg.Wait()

	expect := `           php.txt: PHP(PHP: Hypertext Preprocessor)
    i_am_a_CAT.txt: I am a cat. As yet I have no name.
吾輩は猫である.txt: 吾輩は猫である。名前はまだ無い。
`
	actual := output.String()

	if expect != actual {
		t.Fatal("Incorrect align")
	}
}

type testFileSet struct {
	dir   string
	files []*os.File
}

func newTestFileSet(basenames []string) (*testFileSet, error) {
	dir, err := ioutil.TempDir("", "ninetail")
	if err != nil {
		return nil, err
	}

	files := make([]*os.File, len(basenames))
	for i, name := range basenames {
		filename := filepath.Join(dir, name)
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		files[i] = f
	}

	return &testFileSet{
		dir:   dir,
		files: files,
	}, nil
}

func (tfs *testFileSet) getFilenames() []string {
	filenames := make([]string, len(tfs.files))
	for i, file := range tfs.files {
		filenames[i] = file.Name()
	}
	return filenames
}

func (tfs *testFileSet) writeString(index int, text string) {
	fmt.Fprintf(tfs.files[index], text)
}

func (tfs *testFileSet) removeAll() {
	os.RemoveAll(tfs.dir)
}
