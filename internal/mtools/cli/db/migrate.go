package db

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

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

func printMigrationLog(
	logBuf string,
	namesMap map[string]string,
	err error,
) {
	logLines := strings.Split(logBuf, "\n")
	transformedLines := make([]string, len(logLines))
	applliedLinesMap := make(map[int]struct{})
	for i, line := range logLines {
		after, ok := strings.CutPrefix(line, "Applied: ")
		if !ok {
			transformedLines[i] = line
			continue
		}
		parts := strings.Split(after, " in ")
		after = parts[0]

		migrationFile := strings.TrimSpace(after)
		path, ok := namesMap[migrationFile]
		if !ok {
			transformedLines[i] = line
			continue
		}
		migrationFile = path + "/" + migrationFile

		if !strings.HasPrefix(migrationFile, "/") {
			path, _ := os.Getwd()
			migrationFile = path + "/" + migrationFile
		}

		migrationFile = strings.ReplaceAll(migrationFile, "/./", "/")
		migrationFile = strings.ReplaceAll(migrationFile, "//", "/")
		migrationFile = "file://" + migrationFile

		applliedLinesMap[i] = struct{}{}
		transformedLines[i] = migrationFile
	}

	maxMapKey := -1
	if len(applliedLinesMap) > 0 {
		maxMapKey = slices.Max(slices.Collect(maps.Keys(applliedLinesMap)))
	}
	for i, line := range transformedLines {
		if _, ok := applliedLinesMap[i]; ok {
			if err != nil && maxMapKey == i {
				fmt.Println(color.RedString("Error in Migration: %s", line))
				continue
			}
			fmt.Println(color.GreenString("Applied migration: %s", line))
		} else {
			fmt.Println(line)
		}
	}
}
