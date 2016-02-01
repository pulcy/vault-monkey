package service

import (
	"github.com/hashicorp/vault/api"
)

type VaultService struct {
	vaultClient *api.Client
}

func NewVaultService() (*VaultService, error) {
	// Create a vault client
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, maskAny(err)
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, maskAny(err)
	}

	return &VaultService{
		vaultClient: client,
	}, nil
}
