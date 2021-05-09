package render

import (
	"io"
	"io/ioutil"
	"path/filepath"

	. "github.com/google/go-jsonnet"
	"github.com/gota33/go-config-server/pkg/storage"
)

type Jsonnet struct {
	Importer Importer
}

func (r Jsonnet) Render(entry string, outputType ContentType) (doc string, err error) {
	vm := MakeVM()
	vm.Importer(r.Importer)

	doc, err = vm.EvaluateFile(entry)
	return
}

type RoFsImporter struct {
	Fs storage.ReadonlyFs
}

func (s RoFsImporter) Import(importedFrom, importedPath string) (contents Contents, foundAt string, err error) {
	dir, _ := filepath.Split(importedFrom)
	path := filepath.Join(dir, importedPath)

	var (
		fd   io.ReadCloser
		data []byte
	)
	if fd, err = s.Fs.Open(path); err != nil {
		return
	}
	defer func() { _ = fd.Close() }()

	if data, err = ioutil.ReadAll(fd); err != nil {
		return
	}

	contents = MakeContents(string(data))
	foundAt = path
	return
}
