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
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/giantswarm/retry-go"
	"github.com/hashicorp/vault/api"
)

type TokenConfig struct {
	Policies []string
	Role     string
	Template string
	WrapTTL  string
}

// CreateTokenFile creates a new token with given config and stores it in a file.
func (c *AuthenticatedVaultClient) CreateTokenFile(path string, tokenConfig TokenConfig) error {
	if err := ensureDirectoryOf(path, 0755); err != nil {
		return maskAny(err)
	}
	req := &api.TokenCreateRequest{
		Policies: tokenConfig.Policies,
	}
	if tokenConfig.WrapTTL != "" {
		c.vaultClient.SetWrappingLookupFunc(func(operation, path string) string {
			return tokenConfig.WrapTTL
		})
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
		if tokenConfig.WrapTTL != "" {
			token = secret.WrapInfo.Token
		} else {
			token = secret.Auth.ClientToken
		}
		return nil
	}
	if err := retry.Do(op, retry.RetryChecker(IsVault), retry.MaxTries(3)); err != nil {
		return maskAny(err)
	}

	content := []byte(token)
	if tokenConfig.Template != "" {
		// Put token in template
		t, err := template.New("token").Parse(tokenConfig.Template)
		if err != nil {
			return maskAny(err)
		}
		data := struct {
			Token string
		}{
			Token: token,
		}
		var buffer bytes.Buffer
		if err := t.Execute(&buffer, data); err != nil {
			return maskAny(err)
		}
		content = buffer.Bytes()
	}

	if err := ioutil.WriteFile(path, content, 0400); err != nil {
		return maskAny(err)
	}
	return nil
}
