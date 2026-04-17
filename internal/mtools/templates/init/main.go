package main

import (
	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/config"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"

	"go.uber.org/fx"
)

func main() {
	config.LoadDefaultEnv()

	// DO NOT Remove. It will be edited by the `mtools module create` CLI command.
	modules := []*module.Module{
		cli.NewModule(
			cli.SetConfig(
				cli.ModuleConfig{
					Version: "0.1.0",
					Usage:   "Run project commands",
				},
			),
		),
		logger.NewModule(),
	}

	app := fx.New(
		module.BuildFx(modules...),
		logger.WithLoggerOption(),
		cli.InvokeStartCli(),
	)

	app.Run()
}
