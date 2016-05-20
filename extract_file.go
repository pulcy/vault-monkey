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
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pulcy/vault-monkey/service"
)

var (
	cmdExtractFile = &cobra.Command{
		Use:     "file",
		Short:   "Extract a secret into a file.",
		Example: "--target=<file-path> <secret-path>[#<secret-field>]",
		Run:     cmdExtractFileRun,
	}
)

func init() {
	cmdExtract.AddCommand(cmdExtractFile)
}

func cmdExtractFileRun(cmd *cobra.Command, args []string) {
	// Check arguments
	assertArgIsSet(extractFlags.targetFilePath, "--target")
	if len(args) != 1 {
		Exitf("Provide exactly one argument: <path>[#field]")
	}

	// Parse arguments
	secretPath, secretField, err := parseSecretPath(args[0])
	if err != nil {
		Exitf(err.Error())
	}

	// Login
	c, err := serverLogin()
	if err != nil {
		Exitf("Login failed: %#v", err)
	}

	// Create env file
	secret := service.FileSecret{
		SecretPath:  secretPath,
		SecretField: secretField,
	}
	if err := c.CreateSecretFile(extractFlags.targetFilePath, secret); err != nil {
		Exitf("Secret extraction failed: %v", err)
	}
}

func parseSecretPath(arg string) (string, string, error) {
	pf := strings.Split(arg, "#")
	switch len(pf) {
	case 1:
		return arg, "value", nil
	case 2:
		return pf[0], pf[1], nil
	default:
		return "", "", maskAny(fmt.Errorf("expected '<path>[#field]', got '%s'", arg))
	}
}
