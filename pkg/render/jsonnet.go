package render

import (
	"path/filepath"

	"github.com/GotaX/go-config-server/pkg/storage"
	. "github.com/google/go-jsonnet"
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

type StorageImporter struct {
	Storage storage.Storage
}

func (s StorageImporter) Import(importedFrom, importedPath string) (contents Contents, foundAt string, err error) {
	dir, _ := filepath.Split(importedFrom)
	path := filepath.Join(dir, importedPath)
	data, err := s.Storage.Read(path)

	if err == nil {
		contents = MakeContents(data)
		foundAt = path
	}
	return
}
