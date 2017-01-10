package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// IssueK8sUserCertificate issues a new certificate for a specific service.
func (c *CA) IssueK8sUserCertificate(clusterID, userName, outputDir string) error {
	mountPoint := c.createMountPoint(clusterID, "k8s")
	issuePath := path.Join(mountPoint, "issue", roleOperations)

	// Issue certificate
	data := map[string]interface{}{
		"common_name": userName,
	}
	secret, err := c.vaultClient.Logical().Write(issuePath, data)
	if err != nil {
		return maskAny(err)
	}

	// Write output
	if err := writeData(clusterID, "-cert.pem", outputDir, secret.Data["certificate"]); err != nil {
		return maskAny(err)
	}
	if err := writeData(clusterID, "-key.pem", outputDir, secret.Data["private_key"]); err != nil {
		return maskAny(err)
	}
	if err := writeData(clusterID, "-ca.pem", outputDir, secret.Data["issuing_ca"]); err != nil {
		return maskAny(err)
	}

	return nil
}

func writeData(clusterID, fileSuffix, outputDir string, data interface{}) error {
	os.MkdirAll(outputDir, 0755)
	content, ok := data.(string)
	if !ok {
		return maskAny(fmt.Errorf("Expected data to be a string"))
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, clusterID+fileSuffix), []byte(content), 0600); err != nil {
		return maskAny(err)
	}
	return nil
}
