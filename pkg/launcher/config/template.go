package config

import (
	"fmt"
	"strings"
	"text/template"
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

	sb := &strings.Builder{}
	err = tmpl.Execute(sb, templateFields{os, arch})
	if err != nil {
		return "", fmt.Errorf("Could not execute template: %v", err)
	}

	return sb.String(), nil
}
