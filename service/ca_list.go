package service

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path"
	"sort"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/ryanuber/columnize"
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
	lines := []string{"CommonName | Expiration | Serial"}
	for _, k := range keys {
		key := k.(string)
		if line, err := c.showCertificate(clusterID, key); err != nil {
			return maskAny(err)
		} else {
			lines = append(lines, line)
		}
	}
	sort.Strings(lines[1:])
	fmt.Println(columnize.SimpleFormat(lines))

	return nil
}

func (c *CA) showCertificate(clusterID, serial string) (string, error) {
	mountPoint := c.createMountPoint(clusterID, "k8s")
	certPath := path.Join(mountPoint, "cert", serial)
	secret, err := c.vaultClient.Logical().Read(certPath)
	if err != nil {
		return "", maskAny(err)
	}
	if secret == nil {
		return "", maskAny(fmt.Errorf("No secret returned"))
	}

	certPem, ok := secret.Data["certificate"].(string)
	if !ok {
		return "", maskAny(fmt.Errorf("certificate is not string"))
	}
	block, _ := pem.Decode([]byte(certPem))
	if block == nil {
		return "", maskAny(fmt.Errorf("Failed to parse certificate (PEM) data"))
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", maskAny(err)
	}
	var infoExp string
	if time.Now().After(cert.NotAfter) {
		infoExp = "expired"
	} else {
		infoExp = fmt.Sprintf("expires %s", humanize.Time(cert.NotAfter))
	}
	return fmt.Sprintf("%s | %s | %s", cert.Subject.CommonName, infoExp, serial), nil
}
