package render

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	. "github.com/google/go-jsonnet"
	"github.com/gota33/go-config-server/pkg/storage"
)

type Jsonnet struct {
	Importer Importer
	Data     json.RawMessage
}

func (r Jsonnet) Render(entry string, outputType ContentType) (doc string, err error) {
	vm := MakeVM()
	vm.Importer(r.Importer)

	if len(r.Data) == 0 {
		doc, err = vm.EvaluateFile(entry)
	} else {
		snippet := fmt.Sprintf(`local q = import '%s'; q %s`, entry, r.Data)
		doc, err = vm.EvaluateAnonymousSnippet("snippet.jsonnet", snippet)
	}
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
