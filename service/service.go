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
	"github.com/hashicorp/vault/api"
)

type VaultServiceConfig struct {
	VaultAddr string // URL of the vault
	TokenPath string // Path of a file containing the login token
}

type VaultService struct {
	vaultClient *api.Client
}

func NewVaultService(srvCfg VaultServiceConfig) (*VaultService, error) {
	// Create a vault client
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, maskAny(err)
	}
	if srvCfg.VaultAddr != "" {
		config.Address = srvCfg.VaultAddr
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, maskAny(err)
	}
	if srvCfg.TokenPath != "" {
		token, err := readID(srvCfg.TokenPath)
		if err != nil {
			return nil, maskAny(err)
		}
		client.SetToken(token)
	}

	return &VaultService{
		vaultClient: client,
	}, nil
}
