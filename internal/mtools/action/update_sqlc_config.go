package action

import (
	"context"
	"os"

	"github.com/go-modulus/modulus/errors"
	"gopkg.in/yaml.v3"
)

var ErrSqlcDefinitionFileNotFound = errors.New("project_root/sqlc.definition.yaml file not found")
var ErrSqlcTemplateFileNotFound = errors.New("module_path/storage/sqlc.tmpl.yaml file not found")
var ErrCannotParseSqlcDefinition = errors.New("cannot parse sqlc.definition.yaml file")
var ErrNoSqlcTmpl = errors.New("sqlc.tmpl.yaml file does not exist")
var ErrCannotParseSqlcTmpl = errors.New("cannot parse sqlc.tmpl.yaml file")
var ErrCannotUpdateSqlcConfig = errors.New("cannot update sqlc config")

type UpdateSqlcConfig struct {
}

func NewUpdateSqlcConfig() *UpdateSqlcConfig {
	return &UpdateSqlcConfig{}
}

func (c *UpdateSqlcConfig) Update(ctx context.Context, storagePath string, projPath string) error {
	defFile := projPath + "/sqlc.definition.yaml"
	defContent, err := os.ReadFile(defFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrSqlcDefinitionFileNotFound
		}
		return err
	}

	var def map[string]interface{}

	err = yaml.Unmarshal(defContent, &def)
	if err != nil {
		return errors.WithCause(
			errors.WithHint(
				ErrCannotParseSqlcDefinition,
				"Please check the file "+defFile+" content. It has wrong yaml format.",
			), err,
		)
	}

	sqlcTmplFileName := storagePath + "/sqlc.tmpl.yaml"

	if _, err := os.Stat(sqlcTmplFileName); os.IsNotExist(err) {
		return errors.WithHint(
			ErrNoSqlcTmpl,
			"Please check if the file "+sqlcTmplFileName+" exists.",
		)
	}

	tmplContent, err := os.ReadFile(sqlcTmplFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrSqlcTemplateFileNotFound
		}
		return err
	}

	resContent := defContent
	resContent = append(resContent, []byte("\n\n")...)
	resContent = append(resContent, tmplContent...)

	var tmpl map[string]interface{}

	err = yaml.Unmarshal(resContent, &tmpl)
	if err != nil {
		return errors.WithCause(
			errors.WithHint(
				ErrCannotParseSqlcTmpl,
				"Please check the file"+sqlcTmplFileName+"content. It has wrong yaml format.",
			),
			err,
		)
	}
	for key, val := range def {
		tmpl[key] = val
	}

	sqlcFileName := storagePath + "/sqlc.yaml"
	updateErr := errors.WithHint(
		ErrCannotUpdateSqlcConfig,
		"Some issues occurred when the "+sqlcFileName+" file was being updated.",
	)
	_, err = yaml.Marshal(tmpl)
	if err != nil {
		return errors.WithCause(updateErr, err)
	}

	sqlcContent, err := yaml.Marshal(tmpl["sqlc-tmpl"])
	if err != nil {
		return errors.WithCause(updateErr, err)
	}

	err = os.WriteFile(sqlcFileName, sqlcContent, 0644)
	if err != nil {
		return errors.WithCause(updateErr, err)
	}

	return nil
}
