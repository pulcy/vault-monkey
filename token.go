// Copyright (c) 2017 Pulcy.
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
	"strings"

	"github.com/pulcy/vault-monkey/service"
	"github.com/spf13/cobra"
)

var (
	cmdToken = &cobra.Command{
		Use:   "token",
		Short: "Server commands to operate on tokens",
		Run:   showUsage,
	}

	cmdTokenCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a vault token",
		Run:   cmdTokenCreateRun,
	}

	tokenFlags struct {
		path     string
		policies []string
		role     string
		template string
	}
)

func init() {
	cmdToken.AddCommand(cmdTokenCreate)

	cmdTokenCreate.Flags().StringVar(&tokenFlags.path, "path", "", "Path of the file in which the token will be written")
	cmdTokenCreate.Flags().StringSliceVar(&tokenFlags.policies, "policy", nil, " A list of policies for the token")
	cmdTokenCreate.Flags().StringVar(&tokenFlags.role, "role", "", "If set, the token will be created against the given role")
	cmdTokenCreate.Flags().StringVar(&tokenFlags.template, "template", "", "If set, the token will be wrapped in this Go text template (use {{.Token}})")

	cmdMain.AddCommand(cmdToken)
}

func cmdTokenCreateRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(tokenFlags.path, "path")
	assertArgIsSet(strings.Join(tokenFlags.policies, ","), "policy")

	c, _, err := serverLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := c.CreateTokenFile(tokenFlags.path, service.TokenConfig{
		Policies: tokenFlags.policies,
		Role:     tokenFlags.role,
		Template: tokenFlags.template,
	}); err != nil {
		Exitf("Failed to create token: %v", err)
	}
}
