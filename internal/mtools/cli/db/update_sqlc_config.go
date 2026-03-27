package db

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/mtools/internal/manifesto"
	"github.com/go-modulus/mtools/internal/mtools/action"
	"github.com/urfave/cli/v3"
)

type UpdateSQLCConfig struct {
	action *action.UpdateSqlcConfig
}

func NewUpdateSQLCConfig(action *action.UpdateSqlcConfig) *UpdateSQLCConfig {
	return &UpdateSQLCConfig{
		action: action,
	}
}

func NewUpdateSQLCConfigCommand(updateSqlc *UpdateSQLCConfig) *cli.Command {
	return &cli.Command{
		Name: "update-sqlc-config",
		Usage: `Updates the sqlc config file in all modules of the project.
Example: mtools db update-sqlc-config
`,
		Action: updateSqlc.Invoke,
	}
}

func (c *UpdateSQLCConfig) Invoke(
	ctx context.Context,
	cmd *cli.Command,
) error {
	projPath := cmd.String("proj-path")
	manifest, err := manifesto.LoadLocalManifesto(projPath)
	if err != nil {
		return errors.WithTrace(err)
	}
	for _, md := range manifest.Modules {
		if !md.IsLocalModule {
			continue
		}
		storagePath := md.StoragePath(projPath)
		err := c.action.Update(ctx, storagePath, projPath)
		if err != nil {
			if errors.Is(err, action.ErrNoSqlcTmpl) {
				fmt.Println(
					color.YellowString(
						"No %s/storage/sqlc.tmpl.yaml template file found in the module. Skipping...",
						md.LocalPath,
					),
				)
				continue
			}
			fmt.Println(
				color.RedString(
					"Cannot update %s/storage/sqlc.yaml file for the module %s: %s",
					md.LocalPath,
					md.Name,
					err.Error(),
				),
			)
			continue
		}
		fmt.Println(color.GreenString("%s/storage/sqlc.yaml file updated", md.LocalPath))
	}
	return nil
}
