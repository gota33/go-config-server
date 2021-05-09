package storage

import (
	"context"
	"io"
)

type Provider interface {
	// Provide is used to provide readonly filesystem by namespace
	Provide(ctx context.Context, namespace string) (fs ReadonlyFs, err error)
}

type ReadonlyFs interface {
	io.Closer

	// Open open a readonly file for reading, the associated file descriptor has mode O_RDONLY.
	Open(name string) (ReadonlyFile, error)
}

type ReadonlyFile interface {
	io.ReadCloser
}
