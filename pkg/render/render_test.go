package render

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
)

func TestJsonnet(t *testing.T) {
	for _, c := range cases {
		r := Jsonnet{Importer: makeImporter(c.Files)}
		doc, err := r.Render(c.Entry, c.OutputType)
		if assert.NoError(t, err) {
			assert.Equal(t, c.OutputDoc, compactJson(doc))
		}
	}
}

type TestCase struct {
	Entry      string
	Files      map[string]string
	OutputType ContentType
	OutputDoc  string
}

var cases = []TestCase{
	{
		OutputType: JSON,
		Entry:      "example1.jsonnet",
		OutputDoc:  `{"person1":{"name":"Alice","welcome":"Hello Alice!"},"person2":{"name":"Bob","welcome":"Hello Bob!"}}`,
		Files:      map[string]string{"example1.jsonnet": `/* Edit me! */ {person1:{name:"Alice",welcome:"Hello "+self.name+"!",},person2:self.person1{name:"Bob"},}`},
	},
}

func makeImporter(files map[string]string) jsonnet.Importer {
	data := make(map[string]jsonnet.Contents, len(files))
	for name, content := range files {
		data[name] = jsonnet.MakeContents(content)
	}
	return &jsonnet.MemoryImporter{Data: data}
}

func compactJson(input string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		panic(err)
	}
	return buf.String()
}
