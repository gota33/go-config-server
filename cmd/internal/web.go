package internal

import (
	. "github.com/urfave/cli/v2"
	"go-config-server/pkg/service"
)

const (
	flagHttp       = "http"
	flagRepository = "repository"
	flagUsername   = "username"
	flagPassword   = "password"
)

var CmdWeb = &Command{
	Name:   "web",
	Usage:  "Start config server",
	Action: runWeb,
	Flags: []Flag{
		&StringFlag{
			Name:  flagHttp,
			Value: ":8080",
		},
		&StringFlag{
			Name:     flagRepository,
			Required: true,
			Aliases:  []string{"repo"},
		},
		&StringFlag{
			Name:    flagUsername,
			Aliases: []string{"user"},
		},
		&StringFlag{
			Name:    flagPassword,
			Aliases: []string{"pass"},
		},
	},
}

func runWeb(ctx *Context) (err error) {
	srv := service.New(service.Options{
		HttpAddr: ctx.String(flagHttp),
		URL:      ctx.String(flagRepository),
		Username: ctx.String(flagUsername),
		Password: ctx.String(flagPassword),
	})
	return srv.Run(ctx.Context)
}
