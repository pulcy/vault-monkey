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

type AuthenticatedVaultClient struct {
	log         *logging.Logger
	vaultClient *api.Client
}

func (c *AuthenticatedVaultClient) Cluster() Cluster {
	return Cluster{vaultClient: c.vaultClient}
}

func (c *AuthenticatedVaultClient) Job() Job {
	return Job{vaultClient: c.vaultClient}
}

func (c *AuthenticatedVaultClient) Token() string {
	return c.vaultClient.Token()
}

func (c *AuthenticatedVaultClient) CA() CA {
	return CA{log: c.log, vaultClient: c.vaultClient}
}
