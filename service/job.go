package service

import (
	"fmt"

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

func (vs *VaultService) Job() *Job {
	return &Job{vaultClient: vs.vaultClient}
}

// Create creates the app-id mapping for a job with given id.
func (c *Job) Create(jobId, policyName string) error {
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
	path := fmt.Sprintf("auth/app-id/map/app-id/%s", jobId)
	if _, err := c.vaultClient.Logical().Delete(path); err != nil {
		return maskAny(err)
	}
	// TODO remove all user-id mappings for this job-id (don't see a way how yet)
	// TODO remove all tokens created for this app-id (don't see a way how yet)
	return nil
}

// AddCluster creates the user-id mapping for adding a cluster to a job.
func (c *Job) AddCluster(jobId, clusterId string) error {
	userId := uniuri.NewLen(jobUserIdLen)

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

// RemoveCluster removes the user-id mapping for removing a cluster from the list of clusters allowed to run the job.
func (c *Job) RemoveCluster(jobId, clusterId string) error {
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
