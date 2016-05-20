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

	"github.com/juju/errgo"
)

type GithubLoginData struct {
	GithubToken string
	Mount       string // defaults to "github"
}

// GithubLogin performs a standard Github authentication and initializes the vaultClient with the resulting token.
func (s *VaultService) GithubLogin(data GithubLoginData) (*AuthenticatedVaultClient, error) {
	// Perform login
	vaultClient, address, err := s.newUnsealedClient()
	if err != nil {
		return nil, maskAny(err)
	}
	vaultClient.ClearToken()
	logical := vaultClient.Logical()
	loginData := make(map[string]interface{})
	loginData["token"] = data.GithubToken
	if data.Mount == "" {
		data.Mount = "github"
	}
	path := fmt.Sprintf("auth/%s/login", data.Mount)
	s.log.Debugf("write loginData at %s", address)
	if loginSecret, err := logical.Write(path, loginData); err != nil {
		return nil, maskAny(err)
	} else if loginSecret.Auth == nil {
		return nil, maskAny(errgo.WithCausef(nil, VaultError, "missing authentication in secret response"))
	} else {
		// Use token
		vaultClient.SetToken(loginSecret.Auth.ClientToken)
	}

	// We're done
	return s.newAuthenticatedClient(vaultClient), nil
}
