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
	"context"
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/api/v1"
	retry "github.com/pulcy/vault-monkey/deps/github.com/giantswarm/retry-go"
)

// CreateOrUpdateKubernetesSecret extracts one or more secrets and updates fields in a Kubernetes secret.
func (c *AuthenticatedVaultClient) CreateOrUpdateKubernetesSecret(secretName string, secrets ...EnvSecret) error {
	namespace, err := getKubernetesNamespace()
	if err != nil {
		return maskAny(err)
	}

	client, err := k8s.InClusterClient()
	if err != nil {
		return maskAny(err)
	}

	// Get existing secret or initialize new one
	create := false
	ctx := context.Background()
	secret, err := getKubernetesSecret(ctx, secretName, client, namespace)
	if err != nil {
		create = true
		secret.Metadata = &v1.ObjectMeta{
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
	if err := setKubernetesSecret(ctx, secretName, client, namespace, secret, create); err != nil {
		return maskAny(err)
	}

	return nil
}

// getKubernetesNamespace reads the namespace of the current pod from the well known location.
func getKubernetesNamespace() (string, error) {
	raw, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", maskAny(err)
	}
	return strings.TrimSpace(string(raw)), nil
}

func getKubernetesSecret(ctx context.Context, secretName string, c *k8s.Client, namespace string) (v1.Secret, error) {
	ctx = k8s.NamespaceContext(ctx, namespace)
	s, err := c.CoreV1().GetSecret(ctx, secretName)
	if err != nil {
		return v1.Secret{}, maskAny(err)
	}
	return *s, nil
}

func setKubernetesSecret(ctx context.Context, secretName string, c *k8s.Client, namespace string, secret v1.Secret, create bool) error {
	ctx = k8s.NamespaceContext(ctx, namespace)
	api := c.CoreV1()
	if create {
		if _, err := api.CreateSecret(ctx, &secret); err != nil {
			return maskAny(err)
		}
	} else {
		if _, err := api.UpdateSecret(ctx, &secret); err != nil {
			return maskAny(err)
		}
	}

	return nil
}
