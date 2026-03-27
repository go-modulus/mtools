package db

import (
	"bytes"
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/modulus/errors/errsys"
	"github.com/go-modulus/modulus/errors/errtrace"
	"github.com/go-modulus/mtools/internal/mtools/action"
	"github.com/urfave/cli/v3"
)

var ErrMigrationApplyingError = errsys.New("error applying migration", "Migration cannot be applied.")

type Migrate struct {
	action *action.UpdateSqlcConfig
}

func NewMigrate(
	action *action.UpdateSqlcConfig,
) *Migrate {
	return &Migrate{
		action: action,
	}
}

func NewMigrateCommand(updateSqlc *Migrate) *cli.Command {
	return &cli.Command{
		Name: "migrate",
		Usage: `Migrates all migrations in all modules.
Example: mtools db migrate
Example: mtools db migrate --proj-path=/path/to/project/root
`,
		Action: updateSqlc.Invoke,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "local-manifest",
				Usage:   "Local manifest file related to the project root. Default is modules.json",
				Aliases: []string{"lmf"},
			},
		},
	}
}

func (c *Migrate) Invoke(
	ctx context.Context,
	cmd *cli.Command,
) error {
	projPath := cmd.String("proj-path")
	config, err := newPgxConfig(projPath)
	if err != nil {
		fmt.Println(color.RedString("Cannot load the project config: %s", err.Error()))
		return errtrace.Wrap(err)
	}

	manifest := cmd.String("local-manifest")

	projFs, namesMap, err := commonMigrationFs(projPath, manifest)
	if err != nil {
		return errtrace.Wrap(err)
	}

	dbMate := newDBMate(config, projFs, []string{"migration"})
	var logBuf bytes.Buffer
	dbMate.Log = &logBuf
	err = dbMate.CreateAndMigrate()

	printMigrationLog(logBuf.String(), namesMap, err)

	if err != nil {
		return errtrace.Wrap(errors.WithCause(ErrMigrationApplyingError, err))
	}

	fmt.Println(
		color.GreenString(
			"All migrations are processed.",
		),
	)

	return nil
}
