package command

import (
	"github.com/urfave/cli/v2"
	"leaf-codegen/handler"
	"strings"
)

func (c command) Init() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "init --project <project URL>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			project := strings.ToLower(ctx.String("project"))
			c.log.StandardLogger().Infof("[%s] initializing project...", project)
			return handler.GetHandler().Init(project)
		},
	}
}
