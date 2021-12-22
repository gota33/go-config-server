package internal

import (
	"github.com/gota33/go-config-server/pkg/service"
	. "github.com/urfave/cli/v2"
)

const (
	flagHttp       = "http"
	flagRepository = "repository"
	flagUsername   = "username"
	flagPassword   = "password"
	envPrefix      = "GO_CONFIG_SERVER_"
)

var CmdWeb = &Command{
	Name:   "web",
	Usage:  "Start config server",
	Action: runWeb,
	Flags: []Flag{
		&StringFlag{
			Name:    flagHttp,
			Value:   ":80",
			EnvVars: []string{envPrefix + "HTTP"},
		},
		&StringFlag{
			Name:     flagRepository,
			Required: true,
			Aliases:  []string{"repo"},
			EnvVars:  []string{envPrefix + "REPO"},
		},
		&StringFlag{
			Name:    flagUsername,
			Aliases: []string{"user"},
			EnvVars: []string{envPrefix + "USER"},
		},
		&StringFlag{
			Name:    flagPassword,
			Aliases: []string{"pass"},
			EnvVars: []string{envPrefix + "PASS"},
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
