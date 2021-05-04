package handler

import (
	"context"

	"go-config-server/pkg/render"
	"go-config-server/pkg/storage"
)

type App struct {
	Storage  storage.Storage
	Renderer render.Renderer
}

func (a App) Handle(ctx context.Context, namespace, name string) (doc string, err error) {
	if err = a.Storage.Use(ctx, namespace); err != nil {
		return
	}
	return a.Renderer.Render(name, render.JSON)
}
