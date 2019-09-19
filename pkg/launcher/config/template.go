package config

import (
	"fmt"
	"text/template"

	"github.com/setlog/trivrost/pkg/misc"
)

type templateFields struct {
	OS   string
	Arch string
}

func expandPlaceholders(s string, os string, arch string) (string, error) {
	tmpl, err := template.New("deployment-config").Parse(s)
	if err != nil {
		return "", fmt.Errorf("Could not parse template: %v", err)
	}

	dw := &misc.ByteSliceWriter{}
	err = tmpl.Execute(dw, templateFields{os, arch})
	if err != nil {
		return "", fmt.Errorf("Could not execute template: %v", err)
	}

	return string(dw.Data), nil
}
