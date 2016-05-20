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
	cmdExtractEnv = &cobra.Command{
		Use:     "env",
		Short:   "Extract a secret into an environment file.",
		Example: "--target=<file-path> <environment-key>=<secret-path>[#<secret-field>]...",
		Run:     cmdExtractEnvRun,
	}
)

func init() {
	cmdExtract.AddCommand(cmdExtractEnv)
}

func cmdExtractEnvRun(cmd *cobra.Command, args []string) {
	// Check arguments
	assertArgIsSet(extractFlags.targetFilePath, "--target")
	if len(args) == 0 {
		Exitf("Private at least one argument: <key>=<path>[#field]")
	}

	// Parse arguments
	secrets := []service.EnvSecret{}
	for _, arg := range args {
		secret, err := parseEnvSecret(arg)
		if err != nil {
			Exitf(err.Error())
		}
		secrets = append(secrets, secret)
	}

	// Login
	c, err := serverLogin()
	if err != nil {
		Exitf("Login failed: %#v", err)
	}

	// Create env file
	if err := c.CreateEnvironmentFile(extractFlags.targetFilePath, secrets); err != nil {
		Exitf("Secret extraction failed: %v", err)
	}
}

func parseEnvSecret(arg string) (service.EnvSecret, error) {
	kv := strings.Split(arg, "=")
	if len(kv) != 2 {
		return service.EnvSecret{}, maskAny(fmt.Errorf("expected '<key>=<path>[#field]', got '%s'", arg))
	}
	envKey := kv[0]
	secretPath, secretField, err := parseSecretPath(kv[1])
	if err != nil {
		return service.EnvSecret{}, maskAny(err)
	}
	return service.EnvSecret{
		SecretPath:     secretPath,
		SecretField:    secretField,
		EnvironmentKey: envKey,
	}, nil
}
