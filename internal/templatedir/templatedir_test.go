package templatedir

import (
	"embed"
	"github.com/leaanthony/debme"
	"github.com/matryer/is"
	"io/ioutil"
	"os"
	"testing"
	"testing/fstest"
)

//go:embed testdata/embedded
var embeddedData embed.FS

func TestNew(t *testing.T) {
	is2 := is.New(t)
	// We want an FS further down the embedded data
	fs, err := debme.FS(embeddedData, "testdata/embedded")
	is2.NoErr(err)
	templDir := New(fs)
	targetDir, err := ioutil.TempDir(".", "test_results")
	is2.NoErr(err)
	defer func(path string) {
		err := os.RemoveAll(path)
		is2.NoErr(err)
	}(targetDir)

	templDir.IgnoreFile("ignored.txt")
	templDir.SetTemplateFilters([]string{".filtername", ".tmpl"})

	type TestData struct {
		Name string
	}
	err = templDir.Extract(targetDir, &TestData{
		Name: "Roger Mellie",
	})
	is2.NoErr(err)
	expectedFiles := []string{
		"test.go",
		"subdir/included.txt",
		"subdir/sub.go",
		"custom.txt",
	}

	err = fstest.TestFS(os.DirFS(targetDir), expectedFiles...)
	is2.NoErr(err)
}

func TestBad(t *testing.T) {
	is2 := is.New(t)
	// We want an FS further down the embedded data
	fs, err := debme.FS(embeddedData, "testdata/embedded")
	is2.NoErr(err)
	templDir := New(fs)
	templDir.IgnoreFile("ignored.txt")
	templDir.SetTemplateFilters([]string{".filtername", ".tmpl"})

	type TestData struct {
		Name string
	}
	// Try to extract to non-writable directory
	err = os.Mkdir("/tmp/readonly", 0444)
	is2.NoErr(err)
	err = templDir.Extract("/tmp/readonly", &TestData{
		Name: "Roger Mellie",
	})
	is2.True(err != nil)
	err = os.RemoveAll("/tmp/readonly")
	is2.NoErr(err)
}
