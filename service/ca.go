package service

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/juju/errgo"
	logging "github.com/op/go-logging"
)

const (
	caPolicyPathReadTemplate    = `path "%s" { policy = "read" }`
	caPolicyPathWriteTemplate   = `path "%s" { policy = "write" }`
	roleOperations              = "operations"
	compNameKubeServiceAccounts = "kube-serviceaccounts"
)

type CA struct {
	log         *logging.Logger
	vaultClient *api.Client
}

// createMountPoint creates the mointpoint for the PKI secret backend in the vault, based on the given
// cluster-ID and service name.
func (c *CA) createMountPoint(clusterID, service string) string {
	return path.Join("ca", clusterID, "pki", service)
}

// CreateETCDMembers creates a CA that issues ETCD member certificates.
func (c *CA) CreateETCDMembers(clusterID string, force bool) error {
	mountPoint := c.createMountPoint(clusterID, "etcd")
	if err := c.createRoot(mountPoint, force); err != nil {
		return maskAny(err)
	}
	// Set role
	if err := c.createAnyNameRole(mountPoint, "member"); err != nil {
		return maskAny(err)
	}
	// Create certificate issue policy
	policy, err := c.createIssuePolicy(mountPoint, "member")
	if err != nil {
		return maskAny(err)
	}
	// Create token role
	role := fmt.Sprintf("etcd-%s", clusterID)
	if err := c.createTokenRole(role, []string{policy}); err != nil {
		return maskAny(err)
	}
	// Create & allow job
	if err := c.createJob(clusterID, "etcd", "", policy); err != nil {
		return maskAny(err)
	}

	return nil
}

// CreateK8sAll creates CA's that issues K8S member certificates for all K8S components.
// Each component gets its own CA.
func (c *CA) CreateK8sAll(clusterID, domainName string, force bool) error {
	components := []string{"kubelet", "kube-proxy", "kube-apiserver", "kube-controller-manager", "kube-scheduler", "kube-dns"}
	for _, component := range components {
		if err := c.CreateK8s(clusterID, component, domainName, force); err != nil {
			return maskAny(err)
		}
	}
	// Set role for creating user certificates
	mountPoint := c.createMountPoint(clusterID, "k8s")
	if err := c.createAnyNameRole(mountPoint, roleOperations); err != nil {
		return maskAny(err)
	}
	// Create service-accounts token
	if err := c.createK8sServiceAccountToken(clusterID, force); err != nil {
		return maskAny(err)
	}

	return nil
}

// CreateK8s creates a CA that issues K8S member certificates for the various K8S components.
func (c *CA) CreateK8s(clusterID, component, domainName string, force bool) error {
	mountPoint := c.createMountPoint(clusterID, "k8s")
	if err := c.createRoot(mountPoint, force); err != nil {
		return maskAny(err)
	}
	// Set role
	relPath := path.Join(mountPoint, "roles", component)
	allowedDomains := []string{component}
	if domainName != "" {
		allowedDomains = append(allowedDomains, domainName)
	}
	data := map[string]interface{}{
		"allowed_domains":    strings.Join(allowedDomains, ","),
		"allow_bare_domains": "true",
		"allow_subdomains":   "true",
		"max_ttl":            "720h",
	}
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}
	// Create certificate issue policy
	policy, err := c.createIssuePolicy(mountPoint, component)
	if err != nil {
		return maskAny(err)
	}
	// Create token role
	role := fmt.Sprintf("k8s-%s-%s", clusterID, component)
	if err := c.createTokenRole(role, []string{policy}); err != nil {
		return maskAny(err)
	}
	// Create & allow job
	if err := c.createJob(clusterID, "k8s", component, policy); err != nil {
		return maskAny(err)
	}

	return nil
}

// createK8sServiceAccountToken creates a service-account-token secret used by K8S.
// Roles & policies to access it are also created.
func (c *CA) createK8sServiceAccountToken(clusterID string, force bool) error {
	tokenSecretPath := path.Join("secret", clusterID, "k8s", "token")
	// Create and store secret
	if err := c.createK8sServiceAccountTokenSecret(tokenSecretPath, force); err != nil {
		return maskAny(err)
	}
	// Create read secret policy
	policy, err := c.createReadSecretPolicy(tokenSecretPath, compNameKubeServiceAccounts)
	if err != nil {
		return maskAny(err)
	}
	// Create token role
	role := fmt.Sprintf("k8s-%s-%s", clusterID, compNameKubeServiceAccounts)
	if err := c.createTokenRole(role, []string{policy}); err != nil {
		return maskAny(err)
	}
	// Create & allow job
	if err := c.createJob(clusterID, "k8s", compNameKubeServiceAccounts, policy); err != nil {
		return maskAny(err)
	}

	return nil
}

// createK8sServiceAccountTokenSecret creates the actual token secret (if not exists or forced)
func (c *CA) createK8sServiceAccountTokenSecret(secretPath string, force bool) error {
	secretField := "key"
	if !force {
		// Look if secret already exists
		secret, err := c.vaultClient.Logical().Read(secretPath)
		if err != nil {
			return maskAny(errgo.WithCausef(nil, VaultError, "error reading %s: %s", secretPath, err))
		}
		if secret != nil {
			if _, ok := secret.Data[secretField]; ok {
				// Secret already found
				return nil
			}
		}
	}
	// Create secret
	cmd := exec.Command("openssl", "genrsa", "4096")
	output, err := cmd.Output()
	if err != nil {
		return maskAny(errgo.WithCausef(nil, VaultError, "failed to generate token key at %s: %s", secretPath, err))
	}

	// Write secret
	data := map[string]interface{}{
		secretField: string(output),
	}
	if _, err := c.vaultClient.Logical().Write(secretPath, data); err != nil {
		return maskAny(err)
	}
	return nil
}

// createRoot mounts the PKI backend at the given mountpoint and
// creates the root certificate.
func (c *CA) createRoot(mountPoint string, force bool) error {
	// Check if there is already a PKI mounted at the given mountpoint
	mounts, err := c.vaultClient.Sys().ListMounts()
	if err != nil {
		return maskAny(err)
	}
	if _, found := mounts[mountPoint+"/"]; found {
		// Already mounted
		c.log.Debugf("pki already mounted at %s", mountPoint)
		if !force {
			return nil
		}
	}

	// Mount PKI
	c.log.Debugf("mounting pki at %s", mountPoint)
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
	c.log.Debugf("generating root certificate pki at %s", mountPoint)
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
func (c *CA) createIssuePolicy(mountPoint, role string) (string, error) {
	c.log.Debugf("creating issue policy for %s with role %s", mountPoint, role)
	issuePath := path.Join(mountPoint, "issue", role)
	rules := []string{
		fmt.Sprintf(caPolicyPathWriteTemplate, issuePath),
		fmt.Sprintf(caPolicyPathWriteTemplate, "auth/token/create*"),
	}
	name := path.Join(mountPoint, role)
	policy := strings.Join(rules, "\n")
	if err := c.vaultClient.Sys().PutPolicy(name, policy); err != nil {
		return "", maskAny(err)
	}

	return name, nil
}

// createIssuePolicy creates a mountpoint specific role that allows issueing certificates.
func (c *CA) createReadSecretPolicy(secretPath, role string) (string, error) {
	c.log.Debugf("creating read secret policy for %s with role %s", secretPath, role)
	rules := []string{
		fmt.Sprintf(caPolicyPathReadTemplate, secretPath),
		fmt.Sprintf(caPolicyPathWriteTemplate, "auth/token/create*"),
	}
	name := path.Join(secretPath, role)
	policy := strings.Join(rules, "\n")
	if err := c.vaultClient.Sys().PutPolicy(name, policy); err != nil {
		return "", maskAny(err)
	}

	return name, nil
}

// createTokenRole creates a token role with given name and given allowed policies.
func (c *CA) createTokenRole(role string, policies []string) error {
	relPath := path.Join("auth/token/roles", role)
	c.log.Debugf("creating token role at %s", relPath)
	data := map[string]interface{}{
		"period":           "720h",
		"orphan":           "true",
		"allowed_policies": strings.Join(policies, ","),
	}
	if _, err := c.vaultClient.Logical().Write(relPath, data); err != nil {
		return maskAny(err)
	}

	return nil
}

// createJob creates a job such that vault-monkey can authenticate access it.
func (c *CA) createJob(clusterID, service, component, policyName string) error {
	jobID := fmt.Sprintf("ca-%s-pki-%s", clusterID, service)
	if component != "" {
		jobID = fmt.Sprintf("%s-%s", jobID, component)
	}
	c.log.Debugf("creating job %s with policy %s", jobID, policyName)
	j := Job{c.vaultClient}
	if err := j.Create(jobID, policyName); err != nil {
		return maskAny(err)
	}
	if err := j.AllowCluster(jobID, clusterID); err != nil {
		return maskAny(err)
	}
	return nil
}
