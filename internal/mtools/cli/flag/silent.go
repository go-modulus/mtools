package flag

import (
	"github.com/urfave/cli/v3"
)

func NewSilent(usage string) cli.Flag {
	return &cli.BoolFlag{
		Name:    "silent",
		Usage:   usage,
		Aliases: []string{"s"},
	}
}

func SilentValue(cmd *cli.Command) bool {
	return cmd.Bool("silent")
}
