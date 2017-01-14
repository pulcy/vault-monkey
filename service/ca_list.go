package service

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path"
)

// ListK8sCertificates issues a new certificate for a specific service.
func (c *CA) ListK8sCertificates(clusterID string) error {
	mountPoint := c.createMountPoint(clusterID, "k8s")
	listPath := path.Join(mountPoint, "certs")
	secret, err := c.vaultClient.Logical().List(listPath)
	if err != nil {
		return maskAny(err)
	}
	if secret == nil {
		return maskAny(fmt.Errorf("No secret returned"))
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return maskAny(fmt.Errorf("keys is not an []interface{}"))
	}
	for _, k := range keys {
		key := k.(string)
		if err := c.showCertificate(clusterID, key); err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func (c *CA) showCertificate(clusterID, serial string) error {
	mountPoint := c.createMountPoint(clusterID, "k8s")
	certPath := path.Join(mountPoint, "cert", serial)
	secret, err := c.vaultClient.Logical().Read(certPath)
	if err != nil {
		return maskAny(err)
	}
	if secret == nil {
		return maskAny(fmt.Errorf("No secret returned"))
	}

	certPem, ok := secret.Data["certificate"].(string)
	if !ok {
		return maskAny(fmt.Errorf("certificate is not string"))
	}
	block, _ := pem.Decode([]byte(certPem))
	if block == nil {
		return maskAny(fmt.Errorf("Failed to parse certificate (PEM) data"))
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return maskAny(err)
	}
	fmt.Printf("%s (%s - %s)\n", cert.Subject.CommonName, cert.NotBefore, cert.NotAfter)

	return nil
}
