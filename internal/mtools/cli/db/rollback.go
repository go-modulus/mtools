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

var ErrRollbackApplyingError = errsys.New("error applying rollback", "Migration cannot be rolled back.")

type Rollback struct {
	action *action.UpdateSqlcConfig
}

func NewRollback(
	action *action.UpdateSqlcConfig,
) *Rollback {
	return &Rollback{
		action: action,
	}
}

func NewRollbackCommand(updateSqlc *Rollback) *cli.Command {
	return &cli.Command{
		Name: "rollback",
		Usage: `Rollbacks the last applied migration.
Example: mtools db rollback
Example: mtools db rollback --proj-path=/path/to/project/root
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

func (c *Rollback) Invoke(
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
	err = dbMate.Rollback()
	printMigrationLog(logBuf.String(), namesMap, err)
	if err != nil {
		return errtrace.Wrap(errors.WithCause(ErrRollbackApplyingError, err))
	}

	fmt.Println(
		color.GreenString(
			"The last migration is rolled back.",
		),
	)

	return nil
}
