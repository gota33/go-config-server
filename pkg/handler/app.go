package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/gota33/go-config-server/pkg/render"
	"github.com/gota33/go-config-server/pkg/storage"
)

type App struct {
	Provider    storage.Provider
	NewRenderer func(fs storage.ReadonlyFs, name string) (render.Renderer, error)
}

func (a *App) Handle(ctx context.Context, namespace, name string) (doc string, err error) {
	roFs, err := a.Provider.Provide(ctx, namespace)
	if err != nil {
		return
	}

	defer func() { _ = roFs.Close() }()

	if a.NewRenderer == nil {
		a.NewRenderer = newRenderer
	}

	renderer, err := a.NewRenderer(roFs, name)
	if err != nil {
		return
	}

	return renderer.Render(name, render.JSON)
}

func newRenderer(fs storage.ReadonlyFs, name string) (r render.Renderer, err error) {
	if strings.HasSuffix(name, ".jsonnet") {
		r = render.Jsonnet{Importer: render.RoFsImporter{Fs: fs}}
	} else {
		err = fmt.Errorf("unsupported content type: %q", name)
	}
	return
}
