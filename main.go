package main

import (
	"fmt"
	"os"

	"github.com/ARTI7876/worker-service/cmd"
	msentry "github.com/ARTI7876/worker-service/internal/app/monitor/sentry"
	"github.com/ARTI7876/worker-service/internal/pkg/constant"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    constant.AppName,
		Version: constant.Version,
		Usage:   "Worker Service — потребитель событий order.created",
		Commands: []*cli.Command{
			cmd.WebServer(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "Человеко-читаемый формат логов вместо JSON",
			},
		},
	}

	defer msentry.Flush()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
