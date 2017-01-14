package service

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	logging "github.com/op/go-logging"
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

// SetupDefaults ensures all fields are set to their defaults if they are not yet set
func (config *IssueConfig) SetupDefaults(clusterID string) {
	if config.FileMode == 0 {
		config.FileMode = 0600
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
}

// IssueIsNeeded checks the certificate files that will be created by an issue command.
// If they exist and are valid, an issue is not needed and false will be returned.
func (config *IssueConfig) IssueIsNeeded(log *logging.Logger) bool {
	certPath := filepath.Join(config.OutputDir, config.CertificateFileName)
	keyPath := filepath.Join(config.OutputDir, config.KeyFileName)
	caPath := filepath.Join(config.OutputDir, config.CAFileName)
	return !isValidCertificateSet(log, certPath, keyPath, caPath, config.CommonName)
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
	config.SetupDefaults(clusterID)
	os.MkdirAll(config.OutputDir, 0755)
	certPath := filepath.Join(config.OutputDir, config.CertificateFileName)
	keyPath := filepath.Join(config.OutputDir, config.KeyFileName)
	caPath := filepath.Join(config.OutputDir, config.CAFileName)

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

	// Write output
	mode := os.FileMode(config.FileMode)
	if err := writeData(certPath, mode, secret.Data["certificate"]); err != nil {
		return maskAny(err)
	}
	if err := writeData(keyPath, mode, secret.Data["private_key"]); err != nil {
		return maskAny(err)
	}
	if err := writeData(caPath, mode, secret.Data["issuing_ca"]); err != nil {
		return maskAny(err)
	}

	return nil
}

func writeData(filePath string, fileMode os.FileMode, data interface{}) error {
	content, ok := data.(string)
	if !ok {
		return maskAny(fmt.Errorf("Expected data to be a string"))
	}
	if err := ioutil.WriteFile(filePath, []byte(content), fileMode); err != nil {
		return maskAny(err)
	}
	return nil
}

func isValidCertificateSet(log *logging.Logger, certPath, keyPath, caPath, commonName string) bool {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return false
	}
	if _, err := ioutil.ReadFile(keyPath); err != nil {
		return false
	}
	if _, err := ioutil.ReadFile(caPath); err != nil {
		return false
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		log.Debugf("Failed to parse certificate (PEM) data from %s", certPath)
		return false
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Debugf("Failed to parse certificate (x509) from %s", certPath)
		return false
	}
	if cert.Subject.CommonName != commonName {
		log.Debugf("CommonName differs, wanted '%s', got '%s'", commonName, cert.Subject.CommonName)
		return false
	}
	if time.Now().After(cert.NotAfter) {
		log.Debugf("Certificate expired")
		return false
	}
	ttl := cert.NotAfter.Sub(cert.NotBefore)
	renewTime := cert.NotBefore.Add(ttl / 2)
	if time.Now().After(renewTime) {
		log.Debugf("Certificate passed 50% of TTL")
		return false
	}
	log.Debugf("Certificate valid until %s", humanize.Time(cert.NotAfter))
	return true
}
