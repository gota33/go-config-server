package service

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go-config-server/pkg/handler"
	"go-config-server/pkg/render"
	"go-config-server/pkg/storage"
)

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	endpointHealth = "/healthz"
)

type Options struct {
	HttpAddr string
	URL      string
	Username string
	Password string
}

type Service struct {
	httpAddr string
	app      handler.App
}

func New(opts Options) *Service {
	store := &storage.Git{
		URL: opts.URL,
	}
	if opts.Username != "" {
		store.Auth = &http.BasicAuth{
			Username: opts.Username,
			Password: opts.Password,
		}
	}
	renderer := render.Jsonnet{
		Importer: render.StorageImporter{
			Storage: store,
		},
	}
	return &Service{
		httpAddr: opts.HttpAddr,
		app: handler.App{
			Storage:  store,
			Renderer: renderer,
		},
	}
}

func (srv Service) Handle(c *fiber.Ctx) (err error) {
	var (
		ctx       = c.Context()
		name      = c.Params("+")
		namespace = c.Params("namespace")
	)
	doc, err := srv.app.Handle(ctx, namespace, name)
	if err != nil {
		return
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(doc)
}

func (srv Service) Run(ctx context.Context) (err error) {
	server := fiber.New(fiber.Config{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	registerAccessLogger(server)
	registerHealthHandler(server)
	srv.registerAppHandler(server)

	chErr := make(chan error, 1)

	go func() {
		if err := server.Listen(srv.httpAddr); err != nil {
			chErr <- err
		}
	}()

	select {
	case err = <-chErr:
	case <-ctx.Done():
		err = server.Shutdown()
	}
	return
}

func (srv Service) registerAppHandler(server *fiber.App) {
	server.Get("/:namespace/+", srv.Handle)
}

func registerAccessLogger(server *fiber.App) {
	server.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool { return c.Path() == endpointHealth },
	}))
}

func registerHealthHandler(server *fiber.App) {
	server.All(endpointHealth, func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
