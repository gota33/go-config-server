package service

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gota33/go-config-server/pkg/handler"
	"github.com/gota33/go-config-server/pkg/storage"
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
	store := storage.NewGit(opts.URL)

	if opts.Username != "" {
		store.Auth = &http.BasicAuth{
			Username: opts.Username,
			Password: opts.Password,
		}
	}

	return &Service{
		httpAddr: opts.HttpAddr,
		app:      handler.App{Provider: store},
	}
}

func (srv Service) Handle(c *fiber.Ctx) (err error) {
	ctx := c.Context()
	req := handler.Request{
		Name:      c.Params("+"),
		Namespace: c.Params("namespace"),
	}

	if c.Is(".json") {
		if err = c.BodyParser(&req.Data); err != nil {
			return
		}
	}

	doc, err := srv.app.Handle(ctx, req)
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

	registerHealthHandler(server)
	registerAccessLogger(server)
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
	server.Use(logger.New(logger.Config{}))
}

func registerHealthHandler(server *fiber.App) {
	server.Use(endpointHealth, func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
