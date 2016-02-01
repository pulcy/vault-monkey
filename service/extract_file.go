package service

import (
	"io/ioutil"
)

type FileSecret struct {
	SecretPath  string
	SecretField string
}

// CreateSecretFile extracts one secret and creates a file containing
// the secret value.
func (s *VaultService) CreateSecretFile(path string, secret FileSecret) error {
	if err := ensureDirectoryOf(path, 0755); err != nil {
		return maskAny(err)
	}
	value, err := s.extractSecret(secret.SecretPath, secret.SecretField)
	if err != nil {
		return maskAny(err)
	}
	if err := ioutil.WriteFile(path, []byte(value), 0400); err != nil {
		return maskAny(err)
	}
	return nil
}
