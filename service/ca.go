package service

import (
	"path"
	"strconv"

	"github.com/hashicorp/vault/api"
)

type CA struct {
	vaultClient *api.Client
}

func (c *CA) CreateETCD(mountPoint string) error {
	if err := c.createCore(mountPoint); err != nil {
		return maskAny(err)
	}
	// Set role
	relPath := path.Join(mountPoint, "roles/member")
	data := make(map[string]interface{})
	data["allow_any_name"] = "true"
	data["max_ttl"] = strconv.Itoa(720 * 60 * 60)
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	return nil
}

func (c *CA) CreateK8sKubelet(mountPoint string) error {
	if err := c.createCore(mountPoint); err != nil {
		return maskAny(err)
	}
	// Set role
	relPath := path.Join(mountPoint, "roles/kubelet")
	data := make(map[string]interface{})
	data["allowed_domains"] = "kubelet"
	data["allow_bare_domains"] = "true"
	data["allow_subdomains"] = "false"
	data["max_ttl"] = strconv.Itoa(720 * 60 * 60)
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	return nil
}

func (c *CA) createCore(mountPoint string) error {
	// Mount PKI
	relPath := path.Join("sys/mounts", mountPoint)
	data := make(map[string]interface{})
	data["type"] = "pki"
	data["max_lease_ttl"] = strconv.Itoa(87600 * 60 * 60)
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	// Create root certificate
	relPath = path.Join(mountPoint, "root/generate/internal")
	data = make(map[string]interface{})
	data["common_name"] = mountPoint
	data["ttl"] = strconv.Itoa(87600 * 60 * 60)
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	return nil
}
