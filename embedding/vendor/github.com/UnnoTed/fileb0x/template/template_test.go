package template

import (
	"strings"
	"testing"

	"github.com/UnnoTed/fileb0x/compression"
	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/UnnoTed/fileb0x/updater"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	var err error
	files := make(map[string]*file.File)
	files["test_file.txt"] = &file.File{
		Name: "test_file.txt",
		Path: "static/test_file.txt",
		Data: `[]byte("\x12\x34\x56\x78\x10")`,
	}

	dirs := new(dir.Dir)
	dirs.Insert("static/")

	tp := new(Template)

	err = tp.Set("ayy lmao")
	assert.Error(t, err)
	assert.Equal(t, `Error: Template must be "files" or "file"`, err.Error())

	err = tp.Set("files")
	assert.NoError(t, err)
	assert.Equal(t, "files", tp.name)

	defaultCompression := compression.NewGzip()

	tp.Variables = struct {
		ConfigFile       string
		Now              string
		Pkg              string
		Files            map[string]*file.File
		Spread           bool
		DirList          []string
		Compression      *compression.Options
		Tags             string
		Debug            bool
		Remap            string
		Updater          updater.Config
		ModificationHash string
	}{
		Pkg:              "main",
		Files:            files,
		Spread:           false,
		DirList:          dirs.Clean(),
		Compression:      defaultCompression.Options,
		ModificationHash: "asdasdasd",
	}

	tp.template = "wrong {{.Err pudding"
	tmpl, err := tp.Exec()
	assert.Error(t, err)
	assert.Empty(t, tmpl)

	tp.template = "wrong{{if .Error}} pudding {{end}}"
	tmpl, err = tp.Exec()
	assert.Error(t, err)
	assert.Empty(t, tmpl)

	err = tp.Set("files")
	tmpl, err = tp.Exec()
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpl)

	s := string(tmpl)

	assert.True(t, strings.Contains(s, `var FileStaticTestFileTxt = []byte("\x12\x34\x56\x78\x10")`))
	assert.True(t, strings.Contains(s, `err = FS.Mkdir(CTX, "static/", 0777)`))
	assert.True(t, strings.Contains(s, `f, err = FS.OpenFile(CTX, "static/test_file.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)`))

	// now with spread
	err = tp.Set("file")
	assert.NoError(t, err)
	assert.Equal(t, "file", tp.name)

	defaultCompression = compression.NewGzip()

	tp.Variables = struct {
		ConfigFile   string
		Now          string
		Pkg          string
		Path         string
		Name         string
		Dir          [][]string
		Tags         string
		Remap        string
		Data         string
		Compression  *compression.Options
		Updater      updater.Config
		Modified     string
		OriginalPath string
	}{
		Pkg:          "main",
		Path:         files["test_file.txt"].Path,
		Name:         files["test_file.txt"].Name,
		Dir:          dirs.List,
		Data:         files["test_file.txt"].Data,
		Compression:  defaultCompression.Options,
		Modified:     "11901",
		OriginalPath: "",
	}

	tmpl, err = tp.Exec()
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpl)

	s = string(tmpl)

	assert.True(t, strings.Contains(s, `var FileStaticTestFileTxt = []byte("\x12\x34\x56\x78\x10")`))
	assert.True(t, strings.Contains(s, `f, err := FS.OpenFile(CTX, "static/test_file.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)`))
}
