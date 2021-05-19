package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gota33/go-config-server/pkg/render"
	"github.com/gota33/go-config-server/pkg/storage"
)

type App struct {
	Provider    storage.Provider
	NewRenderer func(fs storage.ReadonlyFs, name string, data json.RawMessage) (render.Renderer, error)
}

type Request struct {
	Namespace string
	Name      string
	Data      json.RawMessage
}

func (a *App) Handle(ctx context.Context, req Request) (doc string, err error) {
	roFs, err := a.Provider.Provide(ctx, req.Namespace)
	if err != nil {
		return
	}

	defer func() { _ = roFs.Close() }()

	if a.NewRenderer == nil {
		a.NewRenderer = newRenderer
	}

	renderer, err := a.NewRenderer(roFs, req.Name, req.Data)
	if err != nil {
		return
	}

	return renderer.Render(req.Name, render.JSON)
}

func newRenderer(fs storage.ReadonlyFs, name string, data json.RawMessage) (r render.Renderer, err error) {
	if strings.HasSuffix(name, ".jsonnet") {
		r = render.Jsonnet{
			Importer: render.RoFsImporter{Fs: fs},
			Data:     data,
		}
	} else {
		err = fmt.Errorf("unsupported content type: %q", name)
	}
	return
}
