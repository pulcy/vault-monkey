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
	"encoding/base64"

	k8s "github.com/YakLabs/k8s-client"
	retry "github.com/pulcy/vault-monkey/deps/github.com/giantswarm/retry-go"
)

// CreateOrUpdateKubernetesSecret extracts one or more secrets and updates fields in a Kubernetes secret.
func (c *AuthenticatedVaultClient) CreateOrUpdateKubernetesSecret(client *K8sClient, secretName string, secrets ...EnvSecret) error {
	namespace, err := getKubernetesNamespace()
	if err != nil {
		return maskAny(err)
	}

	// Get existing secret or initialize new one
	create := false
	secret, err := client.getKubernetesSecret(secretName)
	if err != nil {
		create = true
		secret.ObjectMeta = k8s.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		}
	}
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	// Fetch secrets
	for _, envSec := range secrets {
		var value string
		op := func() error {
			var err error
			value, err = c.extractSecret(envSec.SecretPath, envSec.SecretField)
			if err != nil {
				return maskAny(err)
			}
			return nil
		}
		if err := retry.Do(op, retry.RetryChecker(IsVault), retry.MaxTries(3)); err != nil {
			return maskAny(err)
		}
		encodedValue := base64.StdEncoding.EncodeToString([]byte(value))
		secret.Data[envSec.EnvironmentKey] = []byte(encodedValue)
	}

	// Create/update secret
	if err := client.setKubernetesSecret(secretName, secret, create); err != nil {
		return maskAny(err)
	}

	return nil
}
