package main

import (
	"os"

	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/mtools/internal/mtools"
	"github.com/go-modulus/mtools/internal/mtools/cli/flag"
	cli2 "github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	// current path
	path, _ := os.Getwd()
	cliModule := cli.NewModule(
		cli.OverrideErrorHandler[*mtools.CliErrorHandler],
		cli.SetConfig(
			cli.ModuleConfig{
				Version: "0.1.3",
				Usage:   "This is a CLI tool for the Modulus framework. It helps developer to create a new project, add modules, and manage the project.",
				GlobalFlags: []cli2.Flag{
					flag.NewProjPath(path),
				},
			},
		),
	)

	modules := []*module.Module{
		cliModule,
		logger.NewModule(
			logger.SetConfig(
				logger.ModuleConfig{
					Type:         "console",
					App:          "modulus cli",
					FxEventLevel: zap.WarnLevel.String(),
				},
			),
		),
		mtools.NewModule(),
	}

	app := fx.New(
		module.BuildFx(modules...),
		logger.FxLoggerOption(),
		cli.InvokeStartCli(),
	)

	app.Run()
}
