package manifesto

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/modulus/module"
)

var ErrCannotReadEntries = fmt.Errorf("cannot read entries")

type Entrypoint struct {
	LocalPath string `json:"localPath"`
	Name      string `json:"name"`
}

type LocalManifesto struct {
	Modules []module.Manifesto `json:"modules"`
	Entries []Entrypoint       `json:"entries,omitempty"`
}

func (m *LocalManifesto) ReadFromJSON(data []byte) error {
	return json.Unmarshal(data, &m)
}

func (m *LocalManifesto) WriteToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func (m *LocalManifesto) AddModule(module module.Manifesto) {
	m.Modules = append(m.Modules, module)
}

func (m *LocalManifesto) UpdateModule(module module.Manifesto) {
	for i, mod := range m.Modules {
		if mod.Package == module.Package {
			m.Modules[i] = module
			return
		}
	}
	m.AddModule(module)
}

func (m *LocalManifesto) FindLocalModule(moduleName string) (module.Manifesto, bool) {
	for _, mod := range m.Modules {
		if mod.IsLocalModule && strings.EqualFold(mod.Name, moduleName) {
			return mod, true
		}
	}
	return module.Manifesto{}, false
}

func (m *LocalManifesto) LocalModules() []module.Manifesto {
	res := make([]module.Manifesto, 0)
	for _, mod := range m.Modules {
		if mod.IsLocalModule {
			res = append(res, mod)
		}
	}
	return res
}
func (m *LocalManifesto) SaveAsLocalManifest(projPath string) error {
	data, err := m.WriteToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(projPath+"/modules.json", data, 0644)
}

func NewFromFs(manifestFs fs.FS, filename string) (*LocalManifesto, error) {
	data, err := fs.ReadFile(manifestFs, filename)
	if err != nil {
		return nil, err
	}
	m := &LocalManifesto{}
	err = m.ReadFromJSON(data)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func LoadLocalManifesto(projPath string) (*LocalManifesto, error) {
	entries, err := readEntries(projPath)
	if err != nil {
		return nil, errors.WithCause(ErrCannotReadEntries, err)
	}
	res := LocalManifesto{
		Modules: make([]module.Manifesto, 0),
		Entries: entries,
	}
	if fileExists(projPath + "/modules.json") {
		projFs := os.DirFS(projPath)
		manifest, err := NewFromFs(projFs, "modules.json")
		if err != nil {
			return nil, err
		}
		return manifest, nil
	}
	return &res, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readEntries(projPath string) (entries []Entrypoint, err error) {
	folders, err := os.ReadDir(projPath + "/cmd")
	if err != nil {
		return
	}
	entries = make([]Entrypoint, 0, len(folders))
	for _, entry := range folders {
		if entry.IsDir() {
			entryItem := Entrypoint{
				Name: entry.Name(),
			}
			_, err2 := os.Stat(projPath + "/cmd/" + entry.Name() + "/main.go")
			if os.IsNotExist(err2) {
				continue
			}

			if err2 != nil {
				err = err2
				return
			}
			entryItem.LocalPath = "cmd/" + entry.Name() + "/main.go"
			entries = append(entries, entryItem)
		}
	}

	return
}
