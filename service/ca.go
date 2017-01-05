package service

import (
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
)

const (
	caPolicyTemplate = `path "%s" { policy = "write" }`
)

type CA struct {
	vaultClient *api.Client
}

// CreateETCDMembers creates a CA that issues ETCD member certificates.
func (c *CA) CreateETCDMembers(mountPoint string) error {
	if err := c.createRoot(mountPoint); err != nil {
		return maskAny(err)
	}
	// Set role
	if err := c.createAnyNameRole(mountPoint, "member"); err != nil {
		return maskAny(err)
	}
	// Create certificate issue policy
	if err := c.createIssuePolicy(mountPoint, "member"); err != nil {
		return maskAny(err)
	}

	return nil
}

// CreateK8s creates a CA that issues K8S member certificates for the various K8S components.
func (c *CA) CreateK8s(mountPoint, component string) error {
	if err := c.createRoot(mountPoint); err != nil {
		return maskAny(err)
	}
	// Set role
	relPath := path.Join(mountPoint, "roles", component)
	data := make(map[string]interface{})
	data["allowed_domains"] = component
	data["allow_bare_domains"] = "true"
	data["allow_subdomains"] = "false"
	data["max_ttl"] = "720h"
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	return nil
}

// createRoot mounts the PKI backend at the given mountpoint and
// creates the root certificate.
func (c *CA) createRoot(mountPoint string) error {
	// Mount PKI
	info := &api.MountInput{
		Type:        "pki",
		Description: "CA mount for " + mountPoint,
		Config: api.MountConfigInput{
			MaxLeaseTTL: "87600h",
		},
	}
	if err := c.vaultClient.Sys().Mount(mountPoint, info); err != nil {
		return maskAny(err)
	}

	// Create root certificate
	relPath := path.Join(mountPoint, "root/generate/internal")
	data := make(map[string]interface{})
	data["common_name"] = mountPoint
	data["ttl"] = "87600h"
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	return nil
}

// createAnyNameRole creates a CA role that allows any common name and certificates with a TTL up to 30 days.
func (c *CA) createAnyNameRole(mountPoint, role string) error {
	relPath := path.Join(mountPoint, "roles", role)
	data := make(map[string]interface{})
	data["allow_any_name"] = "true"
	data["max_ttl"] = "720h"
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// createIssuePolicy creates a mountpoint specific role that allows issueing certificates.
func (c *CA) createIssuePolicy(mountPoint, role string) error {
	policyPath := path.Join(mountPoint, "issue", role)
	policy := fmt.Sprintf(caPolicyTemplate, policyPath)
	name := path.Join(mountPoint, role)
	if err := c.vaultClient.Sys().PutPolicy(name, policy); err != nil {
		return maskAny(err)
	}

	return nil
}
