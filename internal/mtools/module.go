package mtools

import (
	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/mtools/internal/mtools/action"
	cmdRoot "github.com/go-modulus/mtools/internal/mtools/cli"
	cmdDb "github.com/go-modulus/mtools/internal/mtools/cli/db"
	cmdModule "github.com/go-modulus/mtools/internal/mtools/cli/module"
	"github.com/go-modulus/pgx"
)

func NewModule() *module.Module {
	return module.NewModule("github.com/go-modulus/mtools").
		AddCliCommands(
			cmdDb.NewDbCommand,
			cmdRoot.NewInitProjectCommand,
			cmdModule.NewModuleCommand,
		).
		AddProviders(
			cmdRoot.NewInitProject,
			cmdModule.NewInstall,
			cmdModule.NewCreate,
			cmdModule.NewAddCli,
			cmdModule.NewAddJsonApi,
			cmdModule.NewInitStorage,
			action.NewInstallStorage,
			action.NewUpdateSqlcConfig,
			cmdDb.NewUpdateSQLCConfig,
			cmdDb.NewAdd,
			cmdDb.NewMigrate,
			cmdDb.NewRollback,
			cmdDb.NewGenerate,
			NewCliErrorHandler,
		).
		AddDependencies(
			logger.NewModule(),
			cli.NewModule(),
			pgx.NewModule(),
		)
}
