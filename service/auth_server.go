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

	"github.com/hashicorp/vault/api"
	"github.com/juju/errgo"
)

const (
	clusterAuthPathPrefix  = "secret/cluster-auth/"
	clusterAuthPathTmpl    = clusterAuthPathPrefix + "%s/job/%s"
	clusterAuthUserIdField = "user-id"
	jobIDEnvKey            = "VAULT_MONKEY_JOB_ID"
	clusterIDEnvKey        = "VAULT_MONKEY_CLUSTER_ID"
	machineIDEnvKey        = "VAULT_MONKEY_MACHINE_ID"
)

type ServerLoginData interface {
	JobID() (string, error)
	ClusterID() (string, error)
	MachineID() (string, error)
}

// ServerLogin performs a 2-step login and initializes the vaultClient with the resulting token.
func (s *VaultService) ServerLogin(data ServerLoginData) (*AuthenticatedVaultClient, error) {
	// Read data
	clusterID, err := data.ClusterID()
	if err != nil {
		return nil, maskAny(err)
	}
	clusterID = strings.ToLower(clusterID)
	machineID, err := data.MachineID()
	if err != nil {
		return nil, maskAny(err)
	}
	machineID = strings.ToLower(machineID)
	jobID, err := data.JobID()
	if err != nil {
		return nil, maskAny(err)
	}
	jobID = strings.ToLower(jobID)

	// Prepare client
	vaultClient, address, err := s.newUnsealedClient()
	if err != nil {
		return nil, maskAny(err)
	}

	// Step 1
	if err := s.serverLoginStep1(vaultClient, clusterID, machineID); err != nil {
		return nil, maskAny(err)
	}

	// Read cluster/job specific user-id
	s.log.Debugf("Fetch cluster+job specific user-id at %s", address)
	logical := vaultClient.Logical()
	userIDPath := fmt.Sprintf(clusterAuthPathTmpl, clusterID, jobID)
	s.log.Debugf("Fetch cluster+job specific user-id from %s", userIDPath)
	userIDSecret, err := logical.Read(userIDPath)
	if err != nil {
		return nil, maskAny(err)
	}
	// Fetch user-id field
	if userIDSecret == nil || userIDSecret.Data == nil {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "userIDSecret == nil at '%s'", userIDPath))
	}
	userID, ok := userIDSecret.Data[clusterAuthUserIdField]
	if !ok {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "missing 'user-id' field at '%s'", userIDPath))
	}

	// Step 2
	client, err := s.serverLoginStep2(vaultClient, jobID, userID)
	if err != nil {
		return nil, maskAny(err)
	}

	return client, nil
}

// serverLoginStep1 performs the first login step using app-id authentication.
func (s *VaultService) serverLoginStep1(vaultClient *api.Client, clusterID, machineID string) error {
	// Perform step 1 login
	var err error
	if s.authMethods.IsEnabled(AuthMethodAppRole) {
		s.log.Debug("Step 1 approle login")
		if err = s.appRoleLogin(vaultClient, clusterID, machineID); err == nil {
			return nil
		}
	}
	if s.authMethods.IsEnabled(AuthMethodAppID) {
		s.log.Debug("Step 1 app-id login")
		if err = s.appIDLogin(vaultClient, clusterID, machineID); err == nil {
			return nil
		}
	}
	if err == nil {
		err = fmt.Errorf("No authentication method left")
	}
	return maskAny(err)
}

// serverLoginStep2 performs a the second step of the 2-step login using the app-id authentication method and
// initializes the vaultClient with the resulting token.
func (s *VaultService) serverLoginStep2(vaultClient *api.Client, jobID string, userID interface{}) (*AuthenticatedVaultClient, error) {
	// Perform step 2 login
	s.log.Debug("Step 2 approle login")
	var err error
	if s.authMethods.IsEnabled(AuthMethodAppRole) {
		if err = s.appRoleLogin(vaultClient, jobID, userID); err == nil {
			return s.newAuthenticatedClient(vaultClient), nil
		}
	}
	if s.authMethods.IsEnabled(AuthMethodAppID) {
		s.log.Debug("Step 2 app-id login")
		if err = s.appIDLogin(vaultClient, jobID, userID); err == nil {
			return s.newAuthenticatedClient(vaultClient), nil
		}
	}
	if err == nil {
		err = fmt.Errorf("No authentication method left")
	}
	return nil, maskAny(err)
}

// appRoleLogin attempts an approle login uisng given roleID & secretID.
// On success, the vaultclient's token is updated with the returned login token.
func (s *VaultService) appRoleLogin(vaultClient *api.Client, roleID string, secretID interface{}) error {
	vaultClient.ClearToken()
	logical := vaultClient.Logical()
	step2Data := make(map[string]interface{})
	step2Data["role_id"] = roleID
	step2Data["secret_id"] = secretID
	if loginSecret, err := logical.Write("auth/approle/login", step2Data); err != nil {
		return maskAny(err)
	} else if loginSecret.Auth == nil {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in secret response"))
	} else {
		// Use step2 token
		vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// We're done
	return nil
}

// appIDLogin attempts an app-ID login uisng given appID & userID.
// On success, the vaultclient's token is updated with the returned login token.
func (s *VaultService) appIDLogin(vaultClient *api.Client, appID string, userID interface{}) error {
	vaultClient.ClearToken()
	logical := vaultClient.Logical()
	//s.log.Debugf("app-id %s", jobID)
	//s.log.Debugf("user-id %s", userId)
	step2Data := make(map[string]interface{})
	step2Data["app_id"] = appID
	step2Data["user_id"] = userID
	if loginSecret, err := logical.Write("auth/app-id/login", step2Data); err != nil {
		return maskAny(err)
	} else if loginSecret.Auth == nil {
		return maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in secret response"))
	} else {
		// Use step2 token
		vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// We're done
	return nil
}

type baseServerLoginData struct {
	next ServerLoginData
}

func (d *baseServerLoginData) JobID() (string, error) {
	if d.next == nil {
		return "", maskAny(fmt.Errorf("JobID not set"))
	}
	return d.next.JobID()
}

func (d *baseServerLoginData) ClusterID() (string, error) {
	if d.next == nil {
		return "", maskAny(fmt.Errorf("ClusterID not set"))
	}
	return d.next.ClusterID()
}

func (d *baseServerLoginData) MachineID() (string, error) {
	if d.next == nil {
		return "", maskAny(fmt.Errorf("MachineID not set"))
	}
	return d.next.MachineID()
}

// NewEnvServerLoginData creates a ServerLoginData that attempts to fetch the data from env variables.
func NewEnvServerLoginData(next ServerLoginData) ServerLoginData {
	return &envServerLoginData{baseServerLoginData{next}}
}

type envServerLoginData struct {
	baseServerLoginData
}

func (d *envServerLoginData) JobID() (string, error) {
	if v := os.Getenv(jobIDEnvKey); v != "" {
		return v, nil
	}
	return d.baseServerLoginData.JobID()
}

func (d *envServerLoginData) ClusterID() (string, error) {
	if v := os.Getenv(clusterIDEnvKey); v != "" {
		return v, nil
	}
	return d.baseServerLoginData.ClusterID()
}

func (d *envServerLoginData) MachineID() (string, error) {
	if v := os.Getenv(machineIDEnvKey); v != "" {
		return v, nil
	}
	return d.baseServerLoginData.MachineID()
}

// NewStaticServerLoginData creates a ServerLoginData that attempts to fetch the data arguments given to this call.
func NewStaticServerLoginData(jobID, clusterID, machineID string, next ServerLoginData) ServerLoginData {
	return &staticServerLoginData{baseServerLoginData{next}, jobID, clusterID, machineID}
}

type staticServerLoginData struct {
	baseServerLoginData
	jobID, clusterID, machineID string
}

func (d *staticServerLoginData) JobID() (string, error) {
	if v := d.jobID; v != "" {
		return v, nil
	}
	return d.baseServerLoginData.JobID()
}

func (d *staticServerLoginData) ClusterID() (string, error) {
	if v := d.clusterID; v != "" {
		return v, nil
	}
	return d.baseServerLoginData.ClusterID()
}

func (d *staticServerLoginData) MachineID() (string, error) {
	if v := d.machineID; v != "" {
		return v, nil
	}
	return d.baseServerLoginData.MachineID()
}

// NewFileSystemServerLoginData creates a ServerLoginData that attempts to fetch the data file files given as arguments to this call.
func NewFileSystemServerLoginData(jobIDPath, clusterIDPath, machineIDPath string, next ServerLoginData) ServerLoginData {
	return &fsServerLoginData{baseServerLoginData{next}, jobIDPath, clusterIDPath, machineIDPath}
}

type fsServerLoginData struct {
	baseServerLoginData
	jobIDPath, clusterIDPath, machineIDPath string
}

func (d *fsServerLoginData) JobID() (string, error) {
	if p := d.jobIDPath; p != "" {
		return readID(p)
	}
	return d.baseServerLoginData.JobID()
}

func (d *fsServerLoginData) ClusterID() (string, error) {
	if p := d.clusterIDPath; p != "" {
		return readID(p)
	}
	return d.baseServerLoginData.ClusterID()
}

func (d *fsServerLoginData) MachineID() (string, error) {
	if p := d.machineIDPath; p != "" {
		return readID(p)
	}
	return d.baseServerLoginData.MachineID()
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
