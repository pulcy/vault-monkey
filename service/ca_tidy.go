package service

import (
	"path"
	"time"
)

// TidyOptions contains custom options for tidy actions.
type TidyOptions struct {
	TidyCertificateStore bool
	TidyRevocationList   bool
	SafetyBuffer         time.Duration
}

// TidyETCDCertificates performs cleanup of the expired ETCD certificates.
func (c *ca) TidyETCDCertificates(clusterID string, options TidyOptions) error {
	if err := c.TidyCertificates(clusterID, "etcd", options); err != nil {
		return maskAny(err)
	}
	return nil
}

// TidyK8sCertificates performs cleanup of the expired kubernetes certificates.
func (c *ca) TidyK8sCertificates(clusterID string, options TidyOptions) error {
	if err := c.TidyCertificates(clusterID, "k8s", options); err != nil {
		return maskAny(err)
	}
	return nil
}

// TidyCertificates performs cleanup of expired certificates for a specific service.
func (c *ca) TidyCertificates(clusterID, service string, options TidyOptions) error {
	mountPoint := c.createMountPoint(clusterID, service)
	path := path.Join(mountPoint, "tidy")
	data := make(map[string]interface{})
	data["tidy_cert_store"] = options.TidyCertificateStore
	data["tidy_revocation_list"] = options.TidyRevocationList
	data["safety_buffer"] = options.SafetyBuffer.String()
	if _, err := c.vaultClient.Logical().Write(path, data); err != nil {
		return maskAny(err)
	}
	return nil
}
