// Copyright (c) 2016 Pulcy.
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
