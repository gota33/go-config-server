package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gota33/go-config-server/pkg/render"
	"github.com/gota33/go-config-server/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	s := &MockStorage{}
	r := &MockRenderer{}
	app := App{Provider: s}
	app.NewRenderer = func(fs storage.ReadonlyFs, name string, data json.RawMessage) (render.Renderer, error) { return r, nil }

	_, _ = app.Handle(context.TODO(), Request{})
	assert.True(t, s.useInvoked)
	assert.True(t, r.renderInvoked)
}

type MockStorage struct {
	useInvoked bool
}

func (m *MockStorage) Provide(ctx context.Context, namespace string) (fs storage.ReadonlyFs, err error) {
	m.useInvoked = true
	return &MockFs{}, nil
}

type MockFs struct{}

func (m MockFs) Close() (_ error) { return }

func (m MockFs) Open(name string) (storage.ReadonlyFile, error) {
	return MockFile{}, nil
}

type MockFile struct{}

func (m MockFile) Read(p []byte) (n int, err error) { return }

func (m MockFile) Close() (err error) { return }

type MockRenderer struct {
	renderInvoked bool
}

func (m *MockRenderer) Render(entry string, outputType render.ContentType) (doc string, err error) {
	m.renderInvoked = true
	return
}
