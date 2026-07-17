package utils_test

import (
	"os"
	"testing"

	"github.com/go-modulus/mtools/internal/mtools/utils"
	"github.com/stretchr/testify/assert"
)

func TestCopyFromTemplates(t *testing.T) {
	t.Run(
		"copy file to tmp dir", func(t *testing.T) {
			err := utils.CopyFromTemplates("init/.env.test", "/tmp/.env.test")

			defer os.Remove("/tmp/.env.test")
			assert.NoError(t, err)
			assert.FileExists(t, "/tmp/.env.test")
		},
	)
}
