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

	"github.com/giantswarm/retry-go"
	"github.com/hashicorp/vault/api"
)

type TokenConfig struct {
	Policies []string
	Role     string
}

// CreateTokenFile creates a new token with given config and stores it in a file.
func (c *AuthenticatedVaultClient) CreateTokenFile(path string, tokenConfig TokenConfig) error {
	if err := ensureDirectoryOf(path, 0755); err != nil {
		return maskAny(err)
	}
	req := &api.TokenCreateRequest{
		Policies: tokenConfig.Policies,
	}
	var token string
	op := func() error {
		var secret *api.Secret
		var err error
		if tokenConfig.Role != "" {
			secret, err = c.vaultClient.Auth().Token().CreateWithRole(req, tokenConfig.Role)
		} else {
			secret, err = c.vaultClient.Auth().Token().Create(req)
		}
		if err != nil {
			return maskAny(err)
		}
		token = secret.Auth.ClientToken
		return nil
	}
	if err := retry.Do(op, retry.RetryChecker(IsVault), retry.MaxTries(3)); err != nil {
		return maskAny(err)
	}
	if err := ioutil.WriteFile(path, []byte(token), 0400); err != nil {
		return maskAny(err)
	}
	return nil
}
