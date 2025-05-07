package jinjarndr

import (
	"github.com/flosch/pongo2"
)

func RenderJinjaTemplate(templateStr string, data map[string]any) (string, error) {
	template, err := pongo2.FromString(templateStr)
	if err != nil {
		return templateStr, err
	}
	rendered, err := template.Execute(pongo2.Context(data))
	if err != nil {
		return templateStr, err
	}
	return rendered, nil
}
