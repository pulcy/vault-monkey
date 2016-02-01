package service

import (
	"github.com/juju/errgo"
)

// extractSecret extracts a secret based on given variables
// Call a login method before calling this method.
func (s *VaultService) extractSecret(secretPath, secretField string) (string, error) {
	if secretPath == "" {
		return "", maskAny(errgo.WithCausef(nil, InvalidArgumentError, "path not set"))
	}
	if secretField == "" {
		return "", maskAny(errgo.WithCausef(nil, InvalidArgumentError, "field not set"))
	}

	// Load secret
	secret, err := s.vaultClient.Logical().Read(secretPath)
	if err != nil {
		return "", maskAny(errgo.WithCausef(nil, VaultError, "error reading %s: %s", secretPath, err))
	}
	if secret == nil {
		return "", maskAny(errgo.WithCausef(nil, VaultError, "no value found at %s", secretPath))
	}

	if value, ok := secret.Data[secretField]; !ok {
		return "", maskAny(errgo.WithCausef(nil, VaultError, "no field '%s' found at %s", secretField, secretPath))
	} else {
		return value.(string), nil
	}
}
