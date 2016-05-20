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

type Cluster struct {
	vaultClient *api.Client
}

// Create creates the app-id mapping for a cluster with given id.
// It also creates and uses a policy for accessing only the jobs within the cluster.
func (c *Cluster) Create(clusterId string) error {
	clusterId = strings.ToLower(clusterId)
	policyName, err := c.createClusterPolicy(clusterId)
	if err != nil {
		return maskAny(err)
	}
	path := fmt.Sprintf("auth/app-id/map/app-id/%s", clusterId)
	data := make(map[string]interface{})
	data["value"] = policyName
	data["display_name"] = clusterId
	if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// Delete removes the app-id mapping for a cluster with given id.
// It also removes the policy for accessing only the jobs within the cluster.
func (c *Cluster) Delete(clusterId string) error {
	clusterId = strings.ToLower(clusterId)
	path := fmt.Sprintf("auth/app-id/map/app-id/%s", clusterId)
	if _, err := c.vaultClient.Logical().Delete(path); err != nil {
		return maskAny(err)
	}
	// TODO remove all user-id mappings for this cluster-id (don't see a way how yet)
	// TODO remove all tokens created for this app-id (don't see a way how yet)
	policyName := fmt.Sprintf(clusterPolicyNameTmpl, clusterId)
	if err := c.vaultClient.Sys().DeletePolicy(policyName); err != nil {
		return maskAny(err)
	}
	return nil
}

// AddMachine creates the user-id mapping for adding a machine to a cluster.
func (c *Cluster) AddMachine(clusterId, machineId, cidrBlock string) error {
	clusterId = strings.ToLower(clusterId)
	machineId = strings.ToLower(machineId)
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
	machineId = strings.ToLower(machineId)
	path := fmt.Sprintf("auth/app-id/map/user-id/%s", machineId)
	if _, err := c.vaultClient.Logical().Delete(path); err != nil {
		return maskAny(err)
	}
	return nil
}

// createClusterPolicy creates and writes a policy into the vault for accessing the
// cluster-auth data of the first step of the server authentication.
// It returns the policy name and any error.
func (c *Cluster) createClusterPolicy(clusterId string) (string, error) {
	clusterId = strings.ToLower(clusterId)
	policy := fmt.Sprintf(clusterPolicyTmpl, clusterAuthPathPrefix+clusterId)
	policyName := fmt.Sprintf(clusterPolicyNameTmpl, clusterId)
	if err := c.vaultClient.Sys().PutPolicy(policyName, policy); err != nil {
		return "", maskAny(err)
	}
	return policyName, nil
}
