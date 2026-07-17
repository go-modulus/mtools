package main

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/mtools/internal/mtools"
	"github.com/go-modulus/mtools/internal/mtools/cli/flag"
	cli2 "github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func moduleVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if bi.Main.Version == "" || bi.Main.Version == "(devel)" {
		return "unknown"
	}
	return strings.TrimPrefix(bi.Main.Version, "v")
}

func main() {
	// current path
	path, _ := os.Getwd()
	cliModule := cli.NewModule(
		cli.OverrideErrorHandler[*mtools.CliErrorHandler],
		cli.SetConfig(
			cli.ModuleConfig{
				Version: moduleVersion(),
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
