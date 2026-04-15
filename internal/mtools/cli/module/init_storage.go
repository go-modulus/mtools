package module

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/mtools/internal/mtools/action"
	"github.com/go-modulus/mtools/internal/mtools/cli/flag"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

type InitStorage struct {
	installStorage *action.InstallStorage
}

func NewInitStorage(installStorage *action.InstallStorage) *InitStorage {
	return &InitStorage{installStorage: installStorage}
}

func NewInitStorageCommand(cmd *InitStorage) *cli.Command {
	return &cli.Command{
		Name: "init-storage",
		Usage: `Initialises the storage feature (SQLc, migrations, queries) for an existing local module.
Example: mtools module init-storage
Example without UI: mtools module init-storage --module=mypckg --silent
Example with custom schema: mtools module init-storage --module=mypckg --schema=myschema
`,
		Action: cmd.Invoke,
		Flags: []cli.Flag{
			flag.NewModule("A name of the local module to initialise storage for"),
			flag.NewSilent("Disable interactive prompts; use defaults for all options"),
			&cli.StringFlag{
				Name:  "schema",
				Usage: "PostgreSQL schema name where tables for this module will be placed",
			},
		},
	}
}

func (c *InitStorage) Invoke(ctx context.Context, cmd *cli.Command) error {
	projPath := cmd.String("proj-path")
	isSilent := cmd.Bool("silent")

	md, err := flag.ModuleValue(cmd)
	if err != nil {
		return err
	}

	cfg := action.StorageConfig{
		Schema:             "public",
		GenerateGraphql:    true,
		GenerateFixture:    true,
		GenerateDataloader: true,
		ProjPath:           projPath,
	}

	schemaFlag := cmd.String("schema")
	if schemaFlag != "" {
		cfg.Schema = schemaFlag
	} else if !isSilent {
		cfg.Schema, err = initStorageAskSchema(cfg.Schema)
		if err != nil {
			return err
		}
	}

	if !isSilent {
		cfg.GenerateGraphql, err = initStorageAskYesNo("Do you want to generate GraphQL files from SQL?")
		if err != nil {
			return err
		}
		cfg.GenerateFixture, err = initStorageAskYesNo("Do you want to generate fixture files from SQL?")
		if err != nil {
			return err
		}
		cfg.GenerateDataloader, err = initStorageAskYesNo("Do you want to generate dataloader files from SQL?")
		if err != nil {
			return err
		}
	}

	err = c.installStorage.Install(ctx, md, cfg)
	if err != nil {
		if errors.Is(err, action.ErrCannotInstallSqlc) {
			fmt.Println(
				color.YellowString(
					"Cannot install the storage feature. Please install SQLc manually: https://docs.sqlc.dev/en/latest/overview/install.html",
				),
			)
			return nil
		}
		if errors.Is(err, action.ErrCannotGenerateFiles) {
			fmt.Println(
				color.YellowString(
					"Cannot generate SQLc DTO files. Please check the errors above. Try to run the command `sqlc generate -f sqlc.yaml` manually`.",
				),
			)
			return nil
		}
		return errors.WithTrace(err)
	}

	fmt.Println(color.GreenString("Storage feature has been successfully initialised for module %s.", md.Name))
	return nil
}

func initStorageAskSchema(defaultSchema string) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Enter a PG schema where you want to place tables for this module: ",
		Default: defaultSchema,
	}
	return prompt.Run()
}

func initStorageAskYesNo(label string) (bool, error) {
	sel := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	_, result, err := sel.Run()
	if err != nil {
		return false, errors.WithCause(
			errors.NewWithHint("cannot ask a question", "Cannot ask a question"), err,
		)
	}
	return result == "Yes", nil
}
