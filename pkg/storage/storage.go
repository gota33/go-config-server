package storage

import "context"

type Storage interface {

	// Use is used for switching between namespaces
	Use(ctx context.Context, namespace string) (err error)

	// Read will return content read from given path
	Read(path string) (content string, err error)
}
