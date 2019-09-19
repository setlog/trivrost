package config

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateDeploymentConfig(expandedDeploymentConfig string) error {
	schemaLoader := gojsonschema.NewStringLoader(deploymentConfigSchema)
	documentLoader := gojsonschema.NewStringLoader(expandedDeploymentConfig)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Could not begin validating the file: %v", err)
	}

	if !result.Valid() {
		errorMessage := ""
		for _, desc := range result.Errors() {
			errorMessage += fmt.Sprintf("\n%s", desc)
		}
		return fmt.Errorf(errorMessage)
	}

	return nil
}
