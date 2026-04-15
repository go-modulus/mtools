package module_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-modulus/mtools/internal/mtools/cli/flag"
	"github.com/go-modulus/mtools/internal/mtools/cli/module"
	"github.com/stretchr/testify/require"
)

const localModulesWithLocalModule = `{
  "name": "Modulus framework modules manifest",
  "version": "1.0.0",
  "description": "List of installed modules for the Modulus framework",
  "modules": [
    {
      "name": "urfave cli",
      "package": "github.com/go-modulus/modulus/cli",
      "description": "Adds ability to create cli applications in the Modulus framework.",
      "install": null,
      "version": "1.0.0"
    },
    {
      "name": "mypckg",
      "package": "testproj/internal/mypckg",
      "description": "A local test module",
      "localPath": "internal/mypckg",
      "isLocalModule": true,
      "version": ""
    }
  ]
}`

func TestInitStorage_Invoke(t *testing.T) {
	t.Run(
		"init storage for an existing local module in silent mode", func(t *testing.T) {
			projDir := "/tmp/testproj_initstorage"
			rb := initProjectWithLocalModule(t, projDir)
			defer rb()

			cmd := module.NewInitStorageCommand(initStorage)
			cmd.Flags = append(cmd.Flags, flag.NewProjPath(projDir))
			err := cmd.Run(
				context.Background(),
				[]string{"init-storage", "--module=mypckg", "--silent", "--proj-path=" + projDir},
			)

			storageDir := fmt.Sprintf("%s/internal/mypckg/storage", projDir)
			_, errStorageDir := os.Stat(storageDir)
			_, errMigrationDir := os.Stat(storageDir + "/migration")
			_, errQueryDir := os.Stat(storageDir + "/query")
			tmplYaml, errTmpl := os.ReadFile(storageDir + "/sqlc.tmpl.yaml")
			defStorageYaml, errDef := os.ReadFile(projDir + "/sqlc.definition.yaml")

			t.Log("When init storage for an existing local module in silent mode")
			t.Log("	The error should be nil")
			require.NoError(t, err)
			t.Log("	The storage directory should be created")
			require.NoError(t, errStorageDir)
			t.Log("	The migration directory should be created")
			require.NoError(t, errMigrationDir)
			t.Log("	The query directory should be created")
			require.NoError(t, errQueryDir)
			t.Log("	The sqlc definition file should be created")
			require.NoError(t, errDef)
			require.Contains(t, string(defStorageYaml), "default-overrides: &default-overrides")
			t.Log("	The sqlc template file should be created")
			require.NoError(t, errTmpl)
			require.Contains(t, string(tmplYaml), "sqlc-tmpl")
		},
	)
}

func initProjectWithLocalModule(t *testing.T, projDir string) func() {
	t.Helper()
	if _, err := os.Stat(projDir); os.IsNotExist(err) {
		err = os.MkdirAll(projDir, 0755)
		if err != nil {
			t.Fatal("Cannot create "+projDir+" dir", err)
		}
		moduleDir := projDir + "/internal/mypckg"
		err = os.MkdirAll(moduleDir, 0755)
		if err != nil {
			t.Fatal("Cannot create "+moduleDir+" dir", err)
		}
		createFile(t, projDir, "modules.json", localModulesWithLocalModule)
		createFile(t, projDir, "go.mod", goModFile)
	}
	return func() {
		_ = os.RemoveAll(projDir)
	}
}
