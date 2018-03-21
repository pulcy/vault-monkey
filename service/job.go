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

	"github.com/dchest/uniuri"
	"github.com/hashicorp/vault/api"
	"github.com/juju/errgo"
)

const (
	jobUserIdLen = 40
)

// Job contains all vault methods to configure secrets for a job.
type Job interface {
	// Create creates the authentication mapping for a job with given id.
	Create(jobID, policyName string) error
	// Delete removes the authentication mapping for a job with given id.
	Delete(jobID string) error
	// AllowCluster creates the user-id mapping for allowing a cluster access to the secrets of a job.
	AllowCluster(jobID, clusterID string) error
	// DenyCluster removes the user-id mapping so the cluster is denied access to the secrets of a job.
	DenyCluster(jobID, clusterID string) error
}

// NewJob creates a new Job manipulator for the given vault client.
func NewJob(vaultClient *api.Client, methods AuthMethod) Job {
	return &job{
		vaultClient: vaultClient,
		methods:     methods,
	}
}

type job struct {
	vaultClient *api.Client
	methods     AuthMethod
}

// Create creates the authentication mapping for a job with given id.
func (c *job) Create(jobID, policyName string) error {
	jobID = strings.ToLower(jobID)
	policyName = strings.ToLower(policyName)

	if c.methods.IsEnabled(AuthMethodAppRole) {
		// Create role
		{
			path := fmt.Sprintf("auth/approle/role/%s", jobID)
			data := make(map[string]interface{})
			data["role_name"] = jobID
			data["bind_secret_id"] = true
			data["policies"] = policyName
			data["secret_id_num_uses"] = 0
			if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
				return maskAny(err)
			}
		}
		// Set role_id
		{
			path := fmt.Sprintf("auth/approle/role/%s/role-id", jobID)
			data := make(map[string]interface{})
			data["role_id"] = jobID
			if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
				return maskAny(err)
			}
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		path := fmt.Sprintf("auth/app-id/map/app-id/%s", jobID)
		data := make(map[string]interface{})
		data["value"] = policyName
		data["display_name"] = jobID
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	return nil
}

// Delete removes the authentication mapping for a job with given id.
func (c *job) Delete(jobID string) error {
	jobID = strings.ToLower(jobID)

	if c.methods.IsEnabled(AuthMethodAppRole) {
		path := fmt.Sprintf("auth/approle/role/%s", jobID)
		if _, err := c.vaultClient.Logical().Delete(path); err != nil {
			return maskAny(err)
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		path := fmt.Sprintf("auth/app-id/map/app-id/%s", jobID)
		if _, err := c.vaultClient.Logical().Delete(path); err != nil {
			return maskAny(err)
		}
		// TODO remove all user-id mappings for this job-id (don't see a way how yet)
		// TODO remove all tokens created for this app-id (don't see a way how yet)
	}
	return nil
}

// AllowCluster creates the user-id mapping for allowing a cluster access to the secrets of a job.
func (c *job) AllowCluster(jobID, clusterID string) error {
	jobID = strings.ToLower(jobID)
	clusterID = strings.ToLower(clusterID)
	userID := strings.ToLower(uniuri.NewLen(jobUserIdLen))

	// Create mapping
	userIDPath := fmt.Sprintf(clusterAuthPathTmpl, clusterID, jobID)
	userIDData := make(map[string]interface{})
	userIDData[clusterAuthUserIdField] = userID
	if _, err := c.vaultClient.Logical().Write(userIDPath, userIDData); err != nil {
		return maskAny(err)
	}

	if c.methods.IsEnabled(AuthMethodAppRole) {
		path := fmt.Sprintf("auth/approle/role/%s/custom-secret-id", jobID)
		data := make(map[string]interface{})
		data["secret_id"] = userID
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		// Map user-id
		path := fmt.Sprintf("auth/app-id/map/user-id/%s", userID)
		data := make(map[string]interface{})
		data["value"] = jobID
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	return nil
}

// DenyCluster removes the user-id mapping so the cluster is denied access to the secrets of a job.
func (c *job) DenyCluster(jobID, clusterID string) error {
	jobID = strings.ToLower(jobID)
	clusterID = strings.ToLower(clusterID)
	// Read the user id
	userIDPath := fmt.Sprintf(clusterAuthPathTmpl, clusterID, jobID)
	userIDSecret, err := c.vaultClient.Logical().Read(userIDPath)
	if err != nil {
		return maskAny(err)
	}

	// Fetch user-id field
	userID, ok := userIDSecret.Data[clusterAuthUserIdField]
	if !ok {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing 'user-id' field at '%s'", userIDPath))
	}

	if c.methods.IsEnabled(AuthMethodAppRole) {
		path := fmt.Sprintf("/auth/approle/role/%s/secret-id/destroy", jobID)
		data := make(map[string]interface{})
		data["secret_id"] = userID
		if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
			return maskAny(err)
		}
	}
	if c.methods.IsEnabled(AuthMethodAppID) {
		// Remove user-id map
		path := fmt.Sprintf("auth/app-id/map/user-id/%s", userID)
		if _, err := c.vaultClient.Logical().Delete(path); err != nil {
			return maskAny(err)
		}
	}
	return nil
}
