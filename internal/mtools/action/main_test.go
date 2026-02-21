package action_test

import (
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/test"
	"github.com/go-modulus/mtools/internal/mtools"
	"github.com/go-modulus/mtools/internal/mtools/action"
	"go.uber.org/fx"
	"testing"
)

var (
	installStorage *action.InstallStorage
)

func TestMain(m *testing.M) {
	test.TestMain(
		m,
		module.BuildFx(mtools.NewModule()),
		fx.Populate(
			&installStorage,
		),
	)
}
