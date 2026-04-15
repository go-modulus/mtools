package module

import "github.com/urfave/cli/v3"

func NewModuleCommand(
	create *Create,
	install *Install,
	addCli *AddCli,
	addJsonApi *AddJsonApi,
	initStorage *InitStorage,
) *cli.Command {
	return &cli.Command{
		Name: "module",
		Usage: `A set of commands for modules manipulations.
Example: mtools module
`,
		Commands: []*cli.Command{
			NewCreateCommand(create),
			NewInstallCommand(install),
			NewAddCliCommand(addCli),
			NewAddJsonApiCommand(addJsonApi),
			NewInitStorageCommand(initStorage),
		},
	}
}
