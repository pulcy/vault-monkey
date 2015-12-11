package main

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// extractSecret extracts a secret based on given variables
func extractSecret(secretPath, secretField string) (string, error) {
	if secretPath == "" {
		return "", maskAny(fmt.Errorf("path not set"))
	}
	if secretField == "" {
		return "", maskAny(fmt.Errorf("field not set"))
	}

	// Create a vault client
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return "", maskAny(err)
	}
	client, err := api.NewClient(config)
	if err != nil {
		return "", maskAny(err)
	}

	// Load secret
	secret, err := client.Logical().Read(secretPath)
	if err != nil {
		return "", maskAny(fmt.Errorf("Error reading %s: %s", secretPath, err))
	}
	if secret == nil {
		return "", maskAny(fmt.Errorf("No value found at %s", secretPath))
	}

	if value, ok := secret.Data[secretField]; !ok {
		return "", maskAny(fmt.Errorf("No field '%s' found at %s", secretField, secretPath))
	} else {
		return value.(string), nil
	}
}
