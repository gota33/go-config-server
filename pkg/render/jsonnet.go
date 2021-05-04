package render

import . "github.com/google/go-jsonnet"

type Jsonnet struct {
	Importer Importer
}

func (r Jsonnet) Render(entry string, outputType ContentType) (doc string, err error) {
	vm := MakeVM()
	vm.Importer(r.Importer)

	doc, err = vm.EvaluateFile(entry)
	return
}
