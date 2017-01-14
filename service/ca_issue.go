package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type IssueConfig struct {
	Role                string
	CommonName          string
	AltNames            []string
	IPSans              []string
	OutputDir           string
	CertificateFileName string
	KeyFileName         string
	CAFileName          string
	FileMode            uint32
}

// IssueETCDCertificate issues a new certificate for a specific service.
func (c *CA) IssueETCDCertificate(clusterID string, config IssueConfig) error {
	if config.Role == "" {
		config.Role = roleMember
	}
	if err := c.IssueCertificate(clusterID, "etcd", config); err != nil {
		return maskAny(err)
	}
	return nil
}

// IssueK8sCertificate issues a new certificate for a specific service.
func (c *CA) IssueK8sCertificate(clusterID string, config IssueConfig) error {
	if config.Role == "" {
		config.Role = roleOperations
	}
	if err := c.IssueCertificate(clusterID, "k8s", config); err != nil {
		return maskAny(err)
	}
	return nil
}

// IssueCertificate issues a new certificate for a specific service.
func (c *CA) IssueCertificate(clusterID, service string, config IssueConfig) error {
	if config.FileMode == 0 {
		config.FileMode = 0600
	}
	mountPoint := c.createMountPoint(clusterID, service)
	issuePath := path.Join(mountPoint, "issue", config.Role)

	// Issue certificate
	data := map[string]interface{}{
		"common_name": config.CommonName,
	}
	if len(config.AltNames) > 0 {
		data["alt_names"] = strings.Join(config.AltNames, ",")
	}
	if len(config.IPSans) > 0 {
		data["ip_sans"] = strings.Join(config.IPSans, ",")
	}
	secret, err := c.vaultClient.Logical().Write(issuePath, data)
	if err != nil {
		return maskAny(err)
	}

	if config.CertificateFileName == "" {
		config.CertificateFileName = clusterID + "-cert.pem"
	}
	if config.KeyFileName == "" {
		config.KeyFileName = clusterID + "-key.pem"
	}
	if config.CAFileName == "" {
		config.CAFileName = clusterID + "-ca.pem"
	}

	// Write output
	mode := os.FileMode(config.FileMode)
	if err := writeData(config.CertificateFileName, config.OutputDir, mode, secret.Data["certificate"]); err != nil {
		return maskAny(err)
	}
	if err := writeData(config.KeyFileName, config.OutputDir, mode, secret.Data["private_key"]); err != nil {
		return maskAny(err)
	}
	if err := writeData(config.CAFileName, config.OutputDir, mode, secret.Data["issuing_ca"]); err != nil {
		return maskAny(err)
	}

	return nil
}

func writeData(fileName, outputDir string, fileMode os.FileMode, data interface{}) error {
	os.MkdirAll(outputDir, 0755)
	content, ok := data.(string)
	if !ok {
		return maskAny(fmt.Errorf("Expected data to be a string"))
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, fileName), []byte(content), fileMode); err != nil {
		return maskAny(err)
	}
	return nil
}
