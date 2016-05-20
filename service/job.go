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

type Job struct {
	vaultClient *api.Client
}

// Create creates the app-id mapping for a job with given id.
func (c *Job) Create(jobId, policyName string) error {
	jobId = strings.ToLower(jobId)
	policyName = strings.ToLower(policyName)
	path := fmt.Sprintf("auth/app-id/map/app-id/%s", jobId)
	data := make(map[string]interface{})
	data["value"] = policyName
	data["display_name"] = jobId
	if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// Delete removes the app-id mapping for a job with given id.
func (c *Job) Delete(jobId string) error {
	jobId = strings.ToLower(jobId)
	path := fmt.Sprintf("auth/app-id/map/app-id/%s", jobId)
	if _, err := c.vaultClient.Logical().Delete(path); err != nil {
		return maskAny(err)
	}
	// TODO remove all user-id mappings for this job-id (don't see a way how yet)
	// TODO remove all tokens created for this app-id (don't see a way how yet)
	return nil
}

// AllowCluster creates the user-id mapping for allowing a cluster access to the secrets of a job.
func (c *Job) AllowCluster(jobId, clusterId string) error {
	jobId = strings.ToLower(jobId)
	clusterId = strings.ToLower(clusterId)
	userId := strings.ToLower(uniuri.NewLen(jobUserIdLen))

	// Create mapping
	userIdPath := fmt.Sprintf(clusterAuthPathTmpl, clusterId, jobId)
	userIdData := make(map[string]interface{})
	userIdData[clusterAuthUserIdField] = userId
	if _, err := c.vaultClient.Logical().Write(userIdPath, userIdData); err != nil {
		return maskAny(err)
	}

	// Map user-id
	path := fmt.Sprintf("auth/app-id/map/user-id/%s", userId)
	data := make(map[string]interface{})
	data["value"] = jobId
	if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// DenyCluster removes the user-id mapping so the cluster is denied access to the secrets of a job.
func (c *Job) DenyCluster(jobId, clusterId string) error {
	jobId = strings.ToLower(jobId)
	clusterId = strings.ToLower(clusterId)
	// Read the user id
	userIdPath := fmt.Sprintf(clusterAuthPathTmpl, clusterId, jobId)
	userIdSecret, err := c.vaultClient.Logical().Read(userIdPath)
	if err != nil {
		return maskAny(err)
	}

	// Fetch user-id field
	userId, ok := userIdSecret.Data[clusterAuthUserIdField]
	if !ok {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing 'user-id' field at '%s'", userIdPath))
	}

	// Remove user-id map
	path := fmt.Sprintf("auth/app-id/map/user-id/%s", userId)
	if _, err := c.vaultClient.Logical().Delete(path); err != nil {
		return maskAny(err)
	}
	return nil
}
