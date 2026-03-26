package module_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-modulus/mtools/internal/mtools/cli/flag"
	"github.com/go-modulus/mtools/internal/mtools/cli/module"
	"github.com/stretchr/testify/require"
)

func createModuleInTmpDir(t *testing.T, projDir string) error {
	cmd := module.NewCreateCommand(createModule)
	cmd.Flags = append(cmd.Flags, flag.NewProjPath(projDir))

	return cmd.Run(
		t.Context(), []string{
			"create", "--package=mypckg", "--name=mypckg", "--path=internal", "--proj-path=" + projDir, "--silent",
		},
	)
}

func TestAddJsonApi_Invoke(t *testing.T) {
	t.Run(
		"create json api", func(t *testing.T) {
			projDir := "/tmp/testproj"
			rb := initProject(t, projDir, goModFile)
			defer rb()

			err := createModuleInTmpDir(t, projDir)
			require.NoError(t, err)

			cmd := module.NewAddJsonApiCommand(addJsonApi)
			cmd.Flags = append(cmd.Flags, flag.NewProjPath(projDir))
			err = cmd.Run(
				t.Context(),
				[]string{
					"add-json-api",
					"--name=HelloWorld",
					"--uri=/mypckg/hello-world",
					"--method=GET",
					"--module=mypckg",
					"--proj-path=" + projDir,
					"--silent",
				},
			)

			apiDir := fmt.Sprintf("%s/internal/mypckg/api", projDir)
			_, errDir := os.Stat(apiDir)
			handlerContent, errCont1 := os.ReadFile(fmt.Sprintf("%s/hello_world.go", apiDir))

			t.Log("When create a new json api handler in the module")
			t.Log("	The error should be nil")
			require.NoError(t, err)
			t.Log("	The api directory should be created in the module")
			require.NoError(t, errDir)
			t.Log("	The api handler file should be created")
			require.NoError(t, errCont1)
			require.Contains(
				t, string(handlerContent), "type HelloWorld struct {",
			)
		},
	)
}
