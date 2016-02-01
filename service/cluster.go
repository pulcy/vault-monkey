package service

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

type Cluster struct {
	vaultClient *api.Client
}

func (vs *VaultService) Cluster() *Cluster {
	return &Cluster{vaultClient: vs.vaultClient}
}

// Create creates the app-id mapping for a cluster with given id and policy.
func (c *Cluster) Create(clusterId, policyName string) error {
	path := fmt.Sprintf("auth/app-id/map/app-id/%s", clusterId)
	data := make(map[string]interface{})
	data["value"] = policyName
	data["display_name"] = clusterId
	if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// AddMachine creates the user-id mapping for adding a machine to a cluster.
func (c *Cluster) AddMachine(clusterId, machineId, cidrBlock string) error {
	path := fmt.Sprintf("auth/app-id/map/user-id/%s", machineId)
	data := make(map[string]interface{})
	data["value"] = clusterId
	if cidrBlock != "" {
		data["cidr_block"] = cidrBlock
	}
	if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// RemoveMachine removes the user-id mapping for removing a machine from a cluster.
func (c *Cluster) RemoveMachine(machineId string) error {
	path := fmt.Sprintf("auth/app-id/map/user-id/%s", machineId)
	if _, err := c.vaultClient.Logical().Delete(path); err != nil {
		return maskAny(err)
	}
	return nil
}
