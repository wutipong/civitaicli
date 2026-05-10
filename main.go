package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/civitaicli/download"
)

func main() {
	godotenv.Load()

	command := cli.Command{
		Name: "Civitai CLI",
		Commands: []*cli.Command{
			download.Command(),
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	command.Run(ctx, os.Args)
}
