package cli_test

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/go-modulus/mtools/internal/mtools/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitProject_Invoke(t *testing.T) {
	t.Run(
		"invoke init project", func(t *testing.T) {
			path := "/tmp/test_project"
			var buf bytes.Buffer
			logger := slog.New(
				slog.NewTextHandler(
					&buf, &slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			c := cli.NewInitProject(logger)

			cmd := cli.NewInitProjectCommand(c)
			err := cmd.Set("path", path)
			require.NoError(t, err, "Setting path flag should not return an error")
			err = cmd.Set("name", "test_project")
			require.NoError(t, err, "Setting name flag should not return an error")

			err = c.Invoke(t.Context(), cmd)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(path)
			_, errDir1 := os.Stat(path)

			contentEnvTest, errEnvTest := os.ReadFile(path + "/.env.test")

			t.Log(buf.String())

			t.Log("When initializing the project")
			assert.NoError(t, err, "The CLI command should not return an error")
			assert.NoError(t, errDir1, "The path should exist")

			assert.NoError(t, errEnvTest, "The .env.test file should exist")
			assert.Contains(
				t,
				string(contentEnvTest),
				"APP_ENV=test",
				"The .env.test file should contain the correct APP_ENV",
			)
		},
	)
}
