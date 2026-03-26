package flag

import (
	"github.com/urfave/cli/v3"
)

func NewProjPath(defaultValue string) cli.Flag {
	return &cli.StringFlag{
		Name:    "proj-path",
		Usage:   "Set the path to the project if it is necessary to run the command from another directory",
		Value:   defaultValue,
		Aliases: []string{"p"},
		Sources: cli.EnvVars("PROJECT_PATH"),
	}
}

func ProjPathValue(cmd *cli.Command) string {
	return cmd.String("proj-path")
}
