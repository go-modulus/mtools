package flag_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-modulus/mtools/internal/mtools/cli/flag"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

const localModulesJson = `{
  "name": "Modulus framework modules registry",
  "description": "List of available modules for the Modulus framework",
  "version": "1.0.0",
  "modules": [
    {
      "name": "captcha",
      "package": "github.com/go-modulus/modulus/captcha",
      "description": "Captcha processor that have to be integrated in auth queries to protect against bots registrations.",
      "install": {
        "envVars": [
          {
            "key": "RECAPTCHA_ENABLED",
            "value": "false",
            "comment": ""
          },
          {
            "key": "RECAPTCHA_V2_SECRET",
            "value": "",
            "comment": ""
          },
          {
            "key": "RECAPTCHA_V3_SECRET",
            "value": "",
            "comment": ""
          },
          {
            "key": "RECAPTCHA_V3_THRESHOLD",
            "value": "0.5",
            "comment": ""
          }
        ],
        "files": [
          {
            "sourceUrl": "https://raw.githubusercontent.com/go-modulus/modulus/refs/heads/main/captcha/install/graphql/captcha.graphql",
            "destFile": "internal/captcha/graphql/captcha.graphql"
          }
        ]
      },
      "version": "1.0.0"
    }
  ]
}`

func TestManifestValue_FetchesFromRegistryURL(t *testing.T) {
	const registryURL = "https://raw.githubusercontent.com/go-modulus/registry/refs/heads/main/modules.json"

	cmd := &cli.Command{
		Flags: []cli.Flag{flag.NewRegistry("")},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			manifest, err := flag.RegistryValue(cmd)

			require.NoError(t, err)
			require.NotNil(t, manifest)
			require.NotEmpty(t, manifest.Modules, "registry manifest should contain at least one module")

			return nil
		},
	}

	err := cmd.Run(context.Background(), []string{"test", "--registry", registryURL})
	require.NoError(t, err)
}

func TestManifestValue_FetchesFromLocalPath(t *testing.T) {

	fn := fmt.Sprintf("%s/%s", os.TempDir(), "fetch-from-local-modules.json")
	err := os.WriteFile(fn, []byte(localModulesJson), 0644)
	defer func() { _ = os.Remove(fn) }()
	require.NoError(t, err)

	cmd := &cli.Command{
		Flags: []cli.Flag{flag.NewRegistry("")},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			manifest, err := flag.RegistryValue(cmd)

			require.NoError(t, err)
			require.NotNil(t, manifest)
			require.Len(t, manifest.Modules, 1, "registry manifest should contain one module from the local file")

			return nil
		},
	}

	err = cmd.Run(context.Background(), []string{"test", "--registry", fn})
	require.NoError(t, err)
}
