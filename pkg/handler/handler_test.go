package handler

import (
	"context"
	"testing"

	"github.com/GotaX/go-config-server/pkg/render"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	s := &MockStorage{}
	r := &MockRender{}
	app := App{Storage: s, Renderer: r}

	_, _ = app.Handle(context.TODO(), "", "")
	assert.True(t, s.useInvoked)
	assert.True(t, r.renderInvoked)
}

type MockStorage struct {
	useInvoked bool
}

func (m *MockStorage) Use(ctx context.Context, namespace string) (err error) {
	m.useInvoked = true
	return
}

func (m *MockStorage) Read(path string) (content string, err error) {
	return
}

type MockRender struct {
	renderInvoked bool
}

func (m *MockRender) Render(entry string, outputType render.ContentType) (doc string, err error) {
	m.renderInvoked = true
	return
}
