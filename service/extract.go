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
	"github.com/juju/errgo"
)

// extractSecret extracts a secret based on given variables
// Call a login method before calling this method.
func (c *AuthenticatedVaultClient) extractSecret(secretPath, secretField string) (string, error) {
	if secretPath == "" {
		return "", maskAny(errgo.WithCausef(nil, InvalidArgumentError, "path not set"))
	}
	if secretField == "" {
		return "", maskAny(errgo.WithCausef(nil, InvalidArgumentError, "field not set"))
	}

	// Load secret
	c.log.Infof("Read %s#%s", secretPath, secretField)
	secret, err := c.vaultClient.Logical().Read(secretPath)
	if err != nil {
		return "", maskAny(errgo.WithCausef(nil, VaultError, "error reading %s: %s", secretPath, err))
	}
	if secret == nil {
		return "", maskAny(errgo.WithCausef(nil, SecretNotFoundError, "no value found at %s", secretPath))
	}

	if value, ok := secret.Data[secretField]; !ok {
		return "", maskAny(errgo.WithCausef(nil, SecretNotFoundError, "no field '%s' found at %s", secretField, secretPath))
	} else {
		return value.(string), nil
	}
}
