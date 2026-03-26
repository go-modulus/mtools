package main

import (
	"fmt"

	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/config"
	"github.com/go-modulus/modulus/module"

	"go.uber.org/fx"
)

func main() {
	fmt.Println("Starting the application...")
	config.LoadDefaultEnv()

	// DO NOT Remove. It will be edited by the `mtools module create` CLI command.
	modules := []*module.Module{
		cli.NewModule(
			cli.InvokeStartCli,
			cli.SetConfig(
				cli.ModuleConfig{
					Version: "0.1.0",
					Usage:   "Run project commands",
				},
			),
		),
	}

	app := fx.New(
		module.BuildFx(modules...),
	)

	app.Run()
}
