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
	"github.com/hashicorp/vault/api"
	"github.com/op/go-logging"
)

// AuthenticatedVaultClient holds a vault client that is already authenticated.
type AuthenticatedVaultClient struct {
	log         *logging.Logger
	vaultClient *api.Client
	authMethods AuthMethod
}

// Cluster returns a helper to configure cluster authentication secrets.
func (c *AuthenticatedVaultClient) Cluster() Cluster {
	return NewCluster(c.vaultClient, c.authMethods)
}

// Job returns a helper to configure job authentication secrets.
func (c *AuthenticatedVaultClient) Job() Job {
	return NewJob(c.vaultClient, c.authMethods)
}

// Token returns the current token of the vault client.
func (c *AuthenticatedVaultClient) Token() string {
	return c.vaultClient.Token()
}

// CA returns a helper to configure certificate authority authentication secrets.
func (c *AuthenticatedVaultClient) CA() CA {
	return NewCA(c.log, c.vaultClient, c.authMethods)
}
