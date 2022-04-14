package main

import (
	"github.com/urfave/cli/v2"
	"leaf-codegen/command"
	"leaf-codegen/logger"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "Leaf Code Generator",
		Usage:       "Supporting leaf framework to initialize project",
		Description: "CLI Leaf code generator",
		UsageText:   "leaf-codegen command [command options] [arguments...]",
		Version:     "v1.0.0",
		Commands: []*cli.Command{
			command.GetCommand().Init(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.GetLogger().StandardLogger().Errorf("Run Error: %+v", err.Error())
	}
}
