// Copyright (c) 2016 Epracom Advies.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		line := fmt.Sprintf("%s=%s", envSec.EnvironmentKey, value)
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
