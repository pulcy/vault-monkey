package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EnvSecret struct {
	SecretPath     string
	SecretField    string
	EnvironmentKey string
}

// CreateEnvironmentFile extracts one or more secrets and creates a key=secretValue
// environment file for them.
func (s *VaultService) CreateEnvironmentFile(path string, secrets []EnvSecret) error {
	if err := ensureDirectoryOf(path, 0755); err != nil {
		return maskAny(err)
	}
	lines := []string{}
	for _, envSec := range secrets {
		value, err := s.extractSecret(envSec.SecretPath, envSec.SecretField)
		if err != nil {
			return maskAny(err)
		}
		line := fmt.Sprintf("%s=%s", envSec.EnvironmentKey, strconv.Quote(value))
		lines = append(lines, line)
	}
	content := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(path, []byte(content), 0400); err != nil {
		return maskAny(err)
	}
	return nil
}

// ensureDirectoryOf creates the directory part of the given file path if needed.
func ensureDirectoryOf(path string, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, perm); err != nil {
		return maskAny(err)
	}
	return nil
}
