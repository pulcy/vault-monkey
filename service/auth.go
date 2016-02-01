package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/juju/errgo"
)

const (
	clusterAuthPathPrefix  = "generic/cluster-auth/"
	clusterAuthPathTmpl    = clusterAuthPathPrefix + "%s/job/%s"
	clusterAuthUserIdField = "user-id"
)

type ServerLoginData struct {
	JobID         string
	ClusterIDPath string
	MachineIDPath string
}

// ServerLogin performs a 2-step login and initializes the vaultClient with the resulting token.
func (s *VaultService) ServerLogin(data ServerLoginData) error {
	// Read data
	clusterID, err := readID(data.ClusterIDPath)
	if err != nil {
		return maskAny(err)
	}
	machineID, err := readID(data.MachineIDPath)
	if err != nil {
		return maskAny(err)
	}

	// Perform step 1 login
	s.vaultClient.ClearToken()
	logical := s.vaultClient.Logical()
	step1Data := make(map[string]interface{})
	step1Data["app_id"] = clusterID
	step1Data["user_id"] = machineID
	if loginSecret, err := logical.Write("auth/app-id/login", step1Data); err != nil {
		return maskAny(err)
	} else if loginSecret.Auth == nil {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in step 1 secret response"))
	} else {
		// Use step1 token
		s.vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// Read cluster/job specific user-id
	userIdPath := fmt.Sprintf(clusterAuthPathTmpl, clusterID, data.JobID)
	userIdSecret, err := logical.Read(userIdPath)
	if err != nil {
		s.vaultClient.ClearToken()
		return maskAny(err)
	}

	// Fetch user-id field
	userId, ok := userIdSecret.Data[clusterAuthUserIdField]
	if !ok {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing 'user-id' field at '%s'", userIdPath))
	}

	// Perform step 2 login
	s.vaultClient.ClearToken()
	step2Data := make(map[string]interface{})
	step2Data["app_id"] = data.JobID
	step2Data["user_id"] = userId
	if loginSecret, err := logical.Write("auth/app-id/login", step2Data); err != nil {
		return maskAny(err)
	} else if loginSecret.Auth == nil {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in step 2 secret response"))
	} else {
		// Use step2 token
		s.vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// We're done
	return nil
}

// readID read an id from a file with given path.
func readID(path string) (string, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", maskAny(errgo.WithCausef(nil, InvalidArgumentError, "%s does not exist", path))
		}
		return "", maskAny(err)
	}
	return strings.TrimSpace(string(raw)), nil
}
