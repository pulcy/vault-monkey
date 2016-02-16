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

package main

import (
	"github.com/pulcy/vault-monkey/service"
)

// adminLogin initialized a VaultServices and tries to perform a administrator login (if needed).
func adminLogin() (*service.VaultService, error) {
	assertArgIsSet(globalFlags.githubToken, "-G")
	// Create service
	vs, err := service.NewVaultService(log, globalFlags.VaultServiceConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	// Login with github (if available)
	if err := vs.GithubLogin(service.GithubLoginData{
		GithubToken: globalFlags.githubToken,
	}); err != nil {
		return nil, maskAny(err)
	}

	return vs, nil
}
