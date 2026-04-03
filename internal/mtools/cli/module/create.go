package module

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"regexp"
	"slices"

	"github.com/fatih/color"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/modulus/errors/errsys"
	"github.com/go-modulus/modulus/errors/errtrace"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/mtools/internal/manifesto"
	"github.com/go-modulus/mtools/internal/mtools/action"
	"github.com/go-modulus/mtools/internal/mtools/files"
	"github.com/go-modulus/mtools/internal/mtools/templates"
	"github.com/go-modulus/mtools/internal/mtools/utils"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

var moduleNameRegexp = regexp.MustCompile(`module\s+([a-zA-Z0-9_\-\/\.]+)+`)
var pckgNameRegexp = regexp.MustCompile(`^[a-z]+[a-z0-9]+`)

var (
	ErrCannotSaveLocalManifest = errsys.New("cannot save a local manifest", "Cannot save a local manifesto")
	ErrCannotCreateDirectory   = errsys.New("cannot create a directory", "Cannot create a directory")
	ErrCannotCreateModuleFile  = errsys.New("cannot create a module file", "Cannot create a module file")
	ErrCannotInstallStorage    = errsys.New("cannot install the storage feature", "Cannot install the storage feature")
	ErrCannotInstallGraphQL    = errsys.New("cannot install the graphql feature", "Cannot install the graphql feature")
)

type features struct {
	storage bool
	graphQL bool
}

type TmplVars struct {
	Module     module.Manifesto
	HasStorage bool
}

type Create struct {
	logger         *slog.Logger
	installStorage *action.InstallStorage
}

func NewCreate(
	logger *slog.Logger,
	installStorage *action.InstallStorage,
) *Create {
	return &Create{
		logger:         logger,
		installStorage: installStorage,
	}
}

func NewCreateCommand(createModule *Create) *cli.Command {
	return &cli.Command{
		Name: "create",
		Usage: `Create a boilerplate of the new module and place its files inside the obtained path.
Adds the chosen module to the project and inits it with copying necessary files.
Example: mtools module create
Example without UI: mtools module create --path=internal/mypckg --package=mypckg --name="My package"
Example filling default values without UI: mtools module create --package=mypckg
`,
		Action: createModule.Invoke,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "package",
				Usage:   "A package name of the module Go file",
				Aliases: []string{"pkg"},
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "A local path to the module",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "A name of the module",
			},
			&cli.BoolFlag{
				Name:    "silent",
				Usage:   "Set the silent mode to disable asking the questions",
				Aliases: []string{"s"},
			},
			&cli.StringSliceFlag{
				Name:  "without",
				Usage: "Set the list of features to install the module without. Available values: storage, graphql",
			},
		},
	}
}

func (c *Create) Invoke(
	ctx context.Context,
	cmd *cli.Command,
) (err error) {
	if !cmd.Bool("silent") {
		utils.PrintLogo()
	}
	projPath := cmd.String("proj-path")

	fmt.Println(color.BlueString("Creating a new module..."))

	manifest, err := manifesto.LoadLocalManifesto(projPath)
	if err != nil {
		return errors.WithTrace(err)
	}

	manifestItem, err := c.getManifestItem(cmd, projPath)
	if err != nil {
		return err
	}

	err = c.saveManifestItem(manifestItem, projPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(projPath+"/"+manifestItem.LocalPath, 0755)
	if err != nil {
		return errtrace.Wrap(
			errors.WithMeta(
				errors.WithCause(ErrCannotCreateDirectory, err),
				"path", manifestItem.LocalPath,
			),
		)
	}

	selectedFeatures := c.getFeatures(cmd)

	if selectedFeatures.storage {
		fmt.Println(color.BlueString("Installing the storage feature..."))
		containsMigrator := slices.ContainsFunc(
			manifest.Modules, func(val module.Manifesto) bool {
				return val.Package == "github.com/go-modulus/pgx/migrator"
			},
		)
		if !containsMigrator {
			fmt.Println(color.YellowString("The 'modulus/pgx/migrator' module is not installed. Please install it to use storage feature."))
		}
		containsPgx := slices.ContainsFunc(
			manifest.Modules, func(val module.Manifesto) bool {
				return val.Package == "github.com/go-modulus/pgx"
			},
		)
		if !containsPgx {
			fmt.Println(color.YellowString("The 'modulus/pgx' module is not installed. Please install it to use storage feature."))
		}
		err = c.installStorageFeature(ctx, cmd, manifestItem, projPath)
		if err != nil {
			return errors.WithTrace(errors.WithCause(ErrCannotInstallStorage, err))
		}
	}

	err = c.addModuleFile(manifestItem, projPath, selectedFeatures)
	if err != nil {
		return err
	}

	err = c.updateEntripoints(projPath, manifestItem)
	if err != nil {
		return err
	}

	fmt.Println(
		color.GreenString("Congratulations! Your module is created."),
	)

	return nil
}

func (c *Create) updateEntripoints(
	projPath string,
	md module.Manifesto,
) error {
	fmt.Println(color.BlueString("Updating entrypoints..."))

	manifest, err := manifesto.LoadLocalManifesto(projPath)
	if err != nil {
		return errors.WithTrace(err)
	}
	for _, entry := range manifest.Entries {
		err = files.AddModuleToEntrypoint(md.Package, projPath+"/"+entry.LocalPath)
		if err != nil {
			fmt.Println(
				color.RedString(
					"Cannot add the module %s to the entrypoint %s: %s. Try to type initialization code manually",
					md.Name,
					entry.LocalPath,
					err.Error(),
				),
			)
			continue
		}
	}

	return nil
}

func (c *Create) installStorageFeature(
	ctx context.Context,
	cmd *cli.Command,
	md module.Manifesto,
	projPath string,
) error {
	cfg := action.StorageConfig{
		Schema:             "public",
		GenerateGraphql:    true,
		GenerateFixture:    true,
		GenerateDataloader: true,
		ProjPath:           projPath,
	}
	if !cmd.Bool("silent") {
		schema, err := c.askSchema(cfg.Schema)
		if err != nil {
			return err
		}
		cfg.Schema = schema

		cfg.GenerateGraphql, err = c.askYesNo("Do you want to generate GraphQL files from SQL?")
		if err != nil {
			return err
		}
		cfg.GenerateFixture, err = c.askYesNo("Do you want to generate fixture files from SQL?")
		if err != nil {
			return err
		}
		cfg.GenerateDataloader, err = c.askYesNo("Do you want to generate dataloader files from SQL?")
		if err != nil {
			return err
		}
	}
	err := c.installStorage.Install(ctx, md, cfg)
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
	return nil
}

func (c *Create) getFeatures(
	cmd *cli.Command,
) (res features) {
	res = features{
		storage: true,
		graphQL: true,
	}
	without := cmd.StringSlice("without")
	type feature struct {
		name        string
		description string
		value       *bool
	}
	items := []feature{
		{
			name: "storage",
			description: "The storage feature allows you to work with PostgreSQL.\n" +
				"It includes migrations and SQLc generated files to call DB queries.",
			value: &res.storage,
		},
		{
			name: "graphql",
			description: "The GraphQL feature allows you to work with GraphQL.\n" +
				"It includes the resolvers and GraphQL schemas compatible with gqlgen.",
			value: &res.graphQL,
		},
	}
	for _, w := range without {
		switch w {
		case "storage":
			res.storage = false
			items = slices.DeleteFunc(
				items, func(val feature) bool {
					return val.name == "storage"
				},
			)
		case "graphql":
			res.graphQL = false
			items = slices.DeleteFunc(
				items, func(val feature) bool {
					return val.name == "graphql"
				},
			)
		}
	}
	if len(items) != 0 && !cmd.Bool("silent") {
		for _, item := range items {
			val, err := c.askYesNo("Do you want to install the " + item.name + " feature?")
			if err != nil {
				return
			}
			*item.value = val
		}

	}
	return
}

func (c *Create) askSchema(defSchema string) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Enter a PG schema where you want to place tables for this module: ",
		Default: defSchema,
	}

	return prompt.Run()
}

func (c *Create) askYesNo(label string) (bool, error) {
	sel := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	_, result, err := sel.Run()
	if err != nil {
		return false, errors.WithCause(errors.NewWithHint("cannot ask a question", "Cannot ask a question"), err)
	}

	return result == "Yes", nil
}

func (c *Create) addModuleFile(
	md module.Manifesto,
	projPath string,
	selectedFeatures features,
) error {
	vars := TmplVars{
		Module:     md,
		HasStorage: selectedFeatures.storage,
	}
	tmpl := template.Must(
		template.New("module.go.tmpl").
			ParseFS(
				templates.TemplateFiles,
				"create_module/module.go.tmpl",
			),
	)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := tmpl.ExecuteTemplate(w, "module.go.tmpl", &vars)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	moduleFIlename := md.ModulePath(projPath) + "/module.go"
	err = os.WriteFile(moduleFIlename, b.Bytes(), 0644)
	if err != nil {
		return errors.WithTrace(
			errors.WithMeta(
				errsys.WithCause(ErrCannotCreateModuleFile, err),
				"filename",
				moduleFIlename,
			),
		)
	}
	return nil
}

func (c *Create) saveManifestItem(manifestItem module.Manifesto, projPath string) (err error) {
	manifest, err := manifesto.LoadLocalManifesto(projPath)
	if err != nil {
		return errors.WithTrace(err)
	}
	for _, item := range manifest.Modules {
		if item.Package == manifestItem.Package {
			return errors.NewWithHint(
				"the module is already installed",
				fmt.Sprintf("The module %s is already installed", item.Name),
			)
		}
	}
	manifest.Modules = append(
		manifest.Modules, manifestItem,
	)
	err = manifest.SaveAsLocalManifest(projPath)
	if err != nil {
		return errors.WithTrace(errors.WithCause(ErrCannotSaveLocalManifest, err))
	}
	return nil
}

func (c *Create) getProjModuleName(projPath string) (string, error) {
	if _, err := os.Stat(projPath + "/go.mod"); os.IsNotExist(err) {
		return "", errors.WithCause(
			errors.NewWithHint(
				"go.mod not found",
				"The "+projPath+"/go.mod file is not found. Try to run the command in the root of the project",
			),
			err,
		)
	}
	content, err := os.ReadFile(projPath + "/go.mod")
	if err != nil {
		return "", errors.WithCause(
			errors.NewWithHint("cannot read a go.mod file", "Cannot read a "+projPath+"/go.mod file"),
			err,
		)
	}

	moduleStr := moduleNameRegexp.FindStringSubmatch(string(content))
	if len(moduleStr) < 2 {
		return "", errors.NewWithHint(
			"cannot find a module name in the go.mod file",
			fmt.Sprintf("Cannot find a module name in the %s/go.mod file in the ", projPath),
		)
	}

	return moduleStr[1], nil
}

func (c *Create) getManifestItem(
	cmd *cli.Command,
	projPath string,
) (
	res module.Manifesto,
	err error,
) {
	isSilent := cmd.Bool("silent")
	pckg := cmd.String("package")
	if pckg == "" {
		if isSilent {
			return module.Manifesto{}, errors.NewWithHint(
				"the package name is not provided",
				"The package name is not provided. Please add the --package flag or remove the --silent=true flag",
			)
		}
		pckg, err = c.askPackage()
		if err != nil {
			return module.Manifesto{}, errors.WithCause(
				errors.New(
					"Cannot get a package name",
				), err,
			)
		}
	} else {
		if !pckgNameRegexp.MatchString(pckg) {
			hint := fmt.Sprintf(
				"The package name %s is not valid. Please use lowercase latin symbols without spaces",
				pckg,
			)
			return module.Manifesto{}, errors.NewWithHint("the package name is not valid", hint)
		}
	}

	name := cmd.String("name")
	if name == "" {
		if !isSilent {
			name, err = c.askName(pckg)
			if err != nil {
				return module.Manifesto{}, errors.WithCause(
					errors.New(
						"Cannot get a name of the module",
					), err,
				)
			}
		} else {
			// If we are in a silence mode, we need to get a name from the package
			name = pckg
		}
	}

	path := cmd.String("path")
	if path == "" {
		if !isSilent {
			path, err = c.askPath(projPath, pckg)
			if err != nil {
				return module.Manifesto{}, errors.WithCause(
					errors.New(
						"Cannot get a path to the module",
					), err,
				)
			}
		} else {
			path = c.getDefaultPath(pckg)
		}
	} else {
		path += "/" + pckg
	}

	projPckg, err := c.getProjModuleName(projPath)
	if err != nil {
		return module.Manifesto{}, err
	}

	res = module.Manifesto{
		Name:          name,
		Package:       projPckg + "/" + path,
		Description:   "",
		Version:       "",
		LocalPath:     path,
		IsLocalModule: true,
	}
	return res, nil
}

func (c *Create) askPath(projPath string, packageName string) (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter a folder starting from the root of a project (" +
			color.BlueString(projPath) + "): ",
	}

	suggestion := "internal"
	prompt.Default = suggestion

	path, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", path, packageName), nil
}

func (c *Create) askName(packageName string) (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter a name of the module: ",
	}

	prompt.Default = packageName

	return prompt.Run()
}

func (c *Create) getDefaultPath(packageName string) string {
	return fmt.Sprintf("internal/%s", packageName)
}

func (c *Create) askPackage() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter a Golang package name of the created module (e.g. user): ",
	}

	var pckg string
	var err error
	for {
		pckg, err = prompt.Run()
		if err != nil {
			return "", err
		}
		if pckg == "" {
			fmt.Println(color.RedString("The package name cannot be empty"))
			continue
		}
		if !pckgNameRegexp.MatchString(pckg) {
			fmt.Println(
				color.RedString(
					"The package name %s is not valid. Please use lowercase latin symbols without spaces",
					pckg,
				),
			)
			continue
		}
		break
	}

	return pckg, nil
}
