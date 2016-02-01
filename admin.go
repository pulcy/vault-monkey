package main

import (
	"./service"
)

// adminLogin initialized a VaultServices and tries to perform a administrator login (if needed).
func adminLogin() (*service.VaultService, error) {
	// Create service
	vs, err := service.NewVaultService(globalFlags.VaultServiceConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	// TODO login if token == ""

	return vs, nil
}
