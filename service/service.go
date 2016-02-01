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
