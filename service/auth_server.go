// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/juju/errgo"
)

const (
	clusterAuthPathPrefix  = "secret/cluster-auth/"
	clusterAuthPathTmpl    = clusterAuthPathPrefix + "%s/job/%s"
	clusterAuthUserIdField = "user-id"
)

type ServerLoginData struct {
	JobID         string
	ClusterIDPath string
	MachineIDPath string
}

// ServerLogin performs a 2-step login and initializes the vaultClient with the resulting token.
func (s *VaultService) ServerLogin(data ServerLoginData) (*AuthenticatedVaultClient, error) {
	// Read data
	clusterID, err := readID(data.ClusterIDPath)
	if err != nil {
		return nil, maskAny(err)
	}
	clusterID = strings.ToLower(clusterID)
	machineID, err := readID(data.MachineIDPath)
	if err != nil {
		return nil, maskAny(err)
	}
	machineID = strings.ToLower(machineID)
	jobID := strings.ToLower(data.JobID)

	// Perform step 1 login
	s.log.Debug("Step 1 login")
	vaultClient, address, err := s.newUnsealedClient()
	if err != nil {
		return nil, maskAny(err)
	}
	vaultClient.ClearToken()
	logical := vaultClient.Logical()
	step1Data := make(map[string]interface{})
	step1Data["app_id"] = clusterID
	step1Data["user_id"] = machineID
	s.log.Debugf("write step1Data to %s", address)
	if loginSecret, err := logical.Write("auth/app-id/login", step1Data); err != nil {
		return nil, maskAny(err)
	} else if loginSecret.Auth == nil {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in step 1 secret response"))
	} else {
		// Use step1 token
		vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// Read cluster/job specific user-id
	s.log.Debugf("Fetch cluster+job specific user-id at %s", address)
	userIdPath := fmt.Sprintf(clusterAuthPathTmpl, clusterID, jobID)
	userIdSecret, err := logical.Read(userIdPath)
	if err != nil {
		return nil, maskAny(err)
	}

	// Fetch user-id field
	if userIdSecret == nil || userIdSecret.Data == nil {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "userIdSecret == nil at '%s'", userIdPath))
	}
	userId, ok := userIdSecret.Data[clusterAuthUserIdField]
	if !ok {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "missing 'user-id' field at '%s'", userIdPath))
	}

	// Perform step 2 login
	s.log.Debug("Step 2 login")
	vaultClient.ClearToken()
	step2Data := make(map[string]interface{})
	step2Data["app_id"] = jobID
	step2Data["user_id"] = userId
	if loginSecret, err := logical.Write("auth/app-id/login", step2Data); err != nil {
		return nil, maskAny(err)
	} else if loginSecret.Auth == nil {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in step 2 secret response"))
	} else {
		// Use step2 token
		vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// We're done
	return s.newAuthenticatedClient(vaultClient), nil
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
