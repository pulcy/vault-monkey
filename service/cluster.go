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
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	clusterPolicyTmpl = `
path "%s/*" {
    policy = "read"
}`
	clusterPolicyNameTmpl = "cluster_auth_%s"
)

// Cluster contains all vault methods to configure secrets for a cluster.
type Cluster interface {
	// Create creates the app-id mapping for a cluster with given id.
	// It also creates and uses a policy for accessing only the jobs within the cluster.
	Create(clusterID string) error
	// Delete removes the app-id mapping for a cluster with given id.
	// It also removes the policy for accessing only the jobs within the cluster.
	Delete(clusterID string) error
	// AddMachine creates the user-id mapping for adding a machine to a cluster.
	AddMachine(clusterID, machineID, cidrBlock string) error
	// RemoveMachine removes the user-id mapping for removing a machine from a cluster.
	RemoveMachine(clusterID, machineID string) error
}

// NewCluster creates a new Cluster manipulator for the given vault client.
func NewCluster(vaultClient *api.Client, methods AuthMethod) Cluster {
	return &cluster{
		vaultClient: vaultClient,
		methods:     methods,
	}
}

type cluster struct {
	vaultClient *api.Client
	methods     AuthMethod
}

// Create creates the app-id mapping for a cluster with given id.
// It also creates and uses a policy for accessing only the jobs within the cluster.
func (c *cluster) Create(clusterID string) error {
	clusterID = strings.ToLower(clusterID)
	policyName, err := c.createClusterPolicy(clusterID)
	if err != nil {
		return maskAny(err)
	}
	if c.methods.IsEnabled(AuthMethodAppRole) {
		// Create role
		{
			path := fmt.Sprintf("auth/approle/role/%s", clusterID)
			data := make(map[string]interface{})
			data["role_name"] = clusterID
			data["bind_secret_id"] = true
			data["policies"] = policyName
			data["secret_id_num_uses"] = 0
			if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
				return maskAny(err)
			}
		}
		// Set role_id
		{
			path := fmt.Sprintf("auth/approle/role/%s/role-id", clusterID)
			data := make(map[string]interface{})
			data["role_id"] = clusterID
			if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
				return maskAny(err)
			}
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		path := fmt.Sprintf("auth/app-id/map/app-id/%s", clusterID)
		data := make(map[string]interface{})
		data["value"] = policyName
		data["display_name"] = clusterID
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	return nil
}

// Delete removes the app-id mapping for a cluster with given id.
// It also removes the policy for accessing only the jobs within the cluster.
func (c *cluster) Delete(clusterID string) error {
	clusterID = strings.ToLower(clusterID)
	if c.methods.IsEnabled(AuthMethodAppRole) {
		path := fmt.Sprintf("auth/approle/role/%s", clusterID)
		if _, err := c.vaultClient.Logical().Delete(path); err != nil {
			return maskAny(err)
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		path := fmt.Sprintf("auth/app-id/map/app-id/%s", clusterID)
		if _, err := c.vaultClient.Logical().Delete(path); err != nil {
			return maskAny(err)
		}
		// TODO remove all user-id mappings for this cluster-id (don't see a way how yet)
		// TODO remove all tokens created for this app-id (don't see a way how yet)
	}
	policyName := fmt.Sprintf(clusterPolicyNameTmpl, clusterID)
	if err := c.vaultClient.Sys().DeletePolicy(policyName); err != nil {
		return maskAny(err)
	}

	return nil
}

// AddMachine creates the user-id mapping for adding a machine to a cluster.
func (c *cluster) AddMachine(clusterID, machineID, cidrBlock string) error {
	clusterID = strings.ToLower(clusterID)
	machineID = strings.ToLower(machineID)
	if c.methods.IsEnabled(AuthMethodAppRole) {
		path := fmt.Sprintf("auth/approle/role/%s/custom-secret-id", clusterID)
		data := make(map[string]interface{})
		data["secret_id"] = machineID
		if cidrBlock != "" {
			data["cidr_list"] = cidrBlock
		}
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		path := fmt.Sprintf("auth/app-id/map/user-id/%s", machineID)
		data := make(map[string]interface{})
		data["value"] = clusterID
		if cidrBlock != "" {
			data["cidr_block"] = cidrBlock
		}
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	return nil
}

// RemoveMachine removes the user-id mapping for removing a machine from a cluster.
func (c *cluster) RemoveMachine(clusterID, machineID string) error {
	clusterID = strings.ToLower(clusterID)
	machineID = strings.ToLower(machineID)
	if c.methods.IsEnabled(AuthMethodAppRole) {
		path := fmt.Sprintf("/auth/approle/role/%s/secret-id/destroy", clusterID)
		data := make(map[string]interface{})
		data["secret_id"] = machineID
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		path := fmt.Sprintf("auth/app-id/map/user-id/%s", machineID)
		if _, err := c.vaultClient.Logical().Delete(path); err != nil {
			return maskAny(err)
		}
	}
	return nil
}

// createClusterPolicy creates and writes a policy into the vault for accessing the
// cluster-auth data of the first step of the server authentication.
// It returns the policy name and any error.
func (c *cluster) createClusterPolicy(clusterID string) (string, error) {
	clusterID = strings.ToLower(clusterID)
	policy := fmt.Sprintf(clusterPolicyTmpl, clusterAuthPathPrefix+clusterID)
	policyName := fmt.Sprintf(clusterPolicyNameTmpl, clusterID)
	if err := c.vaultClient.Sys().PutPolicy(policyName, policy); err != nil {
		return "", maskAny(err)
	}
	return policyName, nil
}
