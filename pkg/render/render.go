package render

const (
	Unknown ContentType = iota
	JSON
	YAML
)

type ContentType int

type Renderer interface {

	// Render returns the rendered document,
	// entry is the entry filename, any file imported by entry file will find from mounted storage.ReadonlyFs,
	// use JSON as the default outputType if outputType is Unknown
	Render(entry string, outputType ContentType) (doc string, err error)
}
