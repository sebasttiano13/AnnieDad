package main

import (
	"context"
	"github.com/sebasttiano13/AnnieDad/internal/application"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
	"github.com/urfave/cli/v3"
	"os"
)

var cfgPath string

func main() {
	cmd := &cli.Command{
		Usage: "AnnieDad server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Server configuration file",
				Destination: &cfgPath,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			return application.Run(ctx, cfgPath)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
}
