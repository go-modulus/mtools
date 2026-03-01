package flag

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/go-modulus/mtools/internal/manifesto"
	"github.com/urfave/cli/v2"
)

func NewManifest(usage string) cli.Flag {
	return &cli.StringFlag{
		Name:        "manifest",
		Usage:       usage,
		DefaultText: "https://raw.githubusercontent.com/go-modulus/registry/refs/heads/main/modules.json",
		Aliases:     []string{"mf"},
		Value:       "https://raw.githubusercontent.com/go-modulus/registry/refs/heads/main/modules.json",
	}
}

func ManifestValue(ctx *cli.Context) (*manifesto.LocalManifesto, error) {
	manifestPath := ctx.String("manifest")
	if manifestPath == "" {
		return manifestFromURL("https://raw.githubusercontent.com/go-modulus/registry/refs/heads/main/modules.json")
	}

	if strings.HasPrefix(manifestPath, "https://") {
		return manifestFromURL(manifestPath)
	}
	if strings.HasPrefix(manifestPath, "http://") {
		return manifestFromURL(manifestPath)
	}

	manifestFs := os.DirFS(filepath.Dir(manifestPath))
	manifestFile := filepath.Base(manifestPath)

	availableModulesManifest, err := manifesto.NewFromFs(manifestFs, manifestFile)
	if err != nil {
		fmt.Println(color.RedString("Cannot read from the manifest file: %s", err.Error()))
		return nil, err
	}
	return availableModulesManifest, nil
}

func manifestFromURL(url string) (*manifesto.LocalManifesto, error) {
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		fmt.Println(color.RedString("Cannot fetch the manifest from URL: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status %s", resp.Status)
		fmt.Println(color.RedString("Cannot fetch the manifest from URL: %s", err.Error()))
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(color.RedString("Cannot read the manifest response body: %s", err.Error()))
		return nil, err
	}

	m := &manifesto.LocalManifesto{}
	if err = m.ReadFromJSON(data); err != nil {
		fmt.Println(color.RedString("Cannot parse the manifest JSON: %s", err.Error()))
		return nil, err
	}
	return m, nil
}
