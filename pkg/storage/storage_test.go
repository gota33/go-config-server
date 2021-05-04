package storage

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	ctx := context.TODO()

	g := &Git{URL: localRepo()}

	if err := g.Use(ctx, "master"); !assert.NoError(t, err) {
		return
	}

	if content, err := g.Read("go.mod"); assert.NoError(t, err) {
		assert.NotEmpty(t, content)
	}
}

func localRepo() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filename, "../../..")
}
