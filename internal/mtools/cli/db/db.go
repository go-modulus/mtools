package db

import (
	"context"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/go-modulus/modulus/config"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/modulus/errors/errsys"
	"github.com/go-modulus/modulus/errors/errtrace"
	"github.com/go-modulus/mtools/internal/manifesto"
	"github.com/go-modulus/pgx"
	"github.com/laher/mergefs"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v3"
)

var ErrCannotLoadManifest = errsys.New("cannot load manifesto", "Cannot load a project manifesto file")

func newDBMate(
	config pgx.ModuleConfig,
	projRootFs fs.FS,
	migrationsDir []string,
) *dbmate.DB {
	u, _ := url.Parse(config.Dsn())
	db := dbmate.New(u)
	db.FS = projRootFs
	db.AutoDumpSchema = false

	db.MigrationsDir = migrationsDir

	return db
}

func newPgxConfig(projPath string) (pgx.ModuleConfig, error) {
	_ = os.Setenv("CONFIG_DIR", projPath)
	config.LoadDefaultEnv()

	cfg := pgx.ModuleConfig{}
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return pgx.ModuleConfig{}, err
	}

	return cfg, nil
}

func commonMigrationFs(projPath string, manifestFile string) (fs.FS, map[string]string, error) {
	projFs := os.DirFS(projPath)
	if manifestFile == "" {
		manifestFile = "modules.json"
	}
	manifest, err := manifesto.NewFromFs(projFs, manifestFile)
	if err != nil {
		return nil, nil, errtrace.Wrap(errors.WithMeta(errors.WithCause(ErrCannotLoadManifest, err), "path", projPath))
	}

	modulesFs := make([]fs.FS, 0)
	fileNamesMap := make(map[string]string)

	for _, md := range manifest.Modules {
		if md.LocalPath == "" {
			continue
		}

		storagePath := md.StoragePath(projPath)
		if _, err := os.Stat(storagePath); os.IsNotExist(err) {
			continue
		}
		localFS := os.DirFS(storagePath)
		modulesFs = append(modulesFs, localFS)
		_ = fs.WalkDir(
			localFS, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				dir := filepath.Dir(path)
				fileNamesMap[filepath.Base(path)] = storagePath + "/" + dir
				return nil
			},
		)
	}

	return mergefs.Merge(modulesFs...), fileNamesMap, nil
}

func NewDbCommand(
	updateSqlc *UpdateSQLCConfig,
	add *Add,
	migrate *Migrate,
	rollback *Rollback,
	generate *Generate,
) *cli.Command {
	return &cli.Command{
		Name: "db",
		Usage: `A set of commands for working with PostgreSQL database in modules.
Example: mtools db
`,
		Commands: []*cli.Command{
			NewUpdateSQLCConfigCommand(updateSqlc),
			NewAddCommand(add),
			NewMigrateCommand(migrate),
			NewRollbackCommand(rollback),
			NewGenerateCommand(generate),
		},
	}
}
