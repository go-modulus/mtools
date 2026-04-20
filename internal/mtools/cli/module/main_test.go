package module_test

import (
	"testing"

	module2 "github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/test"
	"github.com/go-modulus/mtools/internal/mtools"
	"github.com/go-modulus/mtools/internal/mtools/cli/module"
	"go.uber.org/fx"
)

var (
	installModule *module.Install
	createModule  *module.Create
	addJsonApi    *module.AddJsonApi
	initStorage   *module.InitStorage
)

func TestMain(m *testing.M) {
	test.LoadEnv()
	currentModule := mtools.NewModule()
	test.TestMain(
		m,
		module2.BuildFx(currentModule),
		fx.Populate(
			&installModule,
			&createModule,
			&addJsonApi,
			&initStorage,
		),
	)
}
