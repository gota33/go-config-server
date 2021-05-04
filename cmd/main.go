package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"go-config-server/cmd/internal"
)

func main() {
	ctx, cancel := internal.NewAppContext()
	defer cancel()

	app := cli.App{
		Name:     "config-server",
		Commands: []*cli.Command{internal.CmdWeb},
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal("Error: ", err.Error())
	}
}
