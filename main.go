// Copyright (c) 2016 Epracom Advies.
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
	"os"
	"strings"

	"github.com/juju/errgo"
	"github.com/spf13/cobra"

	"git.pulcy.com/pulcy/vault-monkey/service"
)

var (
	projectVersion = "dev"
	projectBuild   = "dev"

	maskAny = errgo.MaskFunc(errgo.Any)
)

var (
	cmdMain = &cobra.Command{
		Use: "vault-monkey",
		Run: showUsage,
	}
	globalFlags struct {
		debug   bool
		verbose bool
		service.VaultServiceConfig
		githubToken string
	}
)

func init() {
	globalFlags.VaultCACert = os.Getenv("VAULT_CACERT")
	globalFlags.VaultCAPath = os.Getenv("VAULT_CAPATH")
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "D", false, "Print debug output")
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
	cmdMain.PersistentFlags().StringVar(&globalFlags.VaultAddr, "vault-addr", "", "URL of the vault (defaults to VAULT_ADDR environment variable)")
	cmdMain.PersistentFlags().StringVar(&globalFlags.VaultCACert, "vault-cacert", globalFlags.VaultCACert, "Path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate")
	cmdMain.PersistentFlags().StringVar(&globalFlags.VaultCAPath, "vault-capath", globalFlags.VaultCAPath, "Path to a directory of PEM-encoded CA cert files to verify the Vault server SSL certificate")
	cmdMain.PersistentFlags().StringVar(&globalFlags.TokenPath, "token-path", "", "Path of a file containing your vault token (token defaults to VAULT_TOKEN environment variable)")
	cmdMain.PersistentFlags().StringVarP(&globalFlags.githubToken, "github-token", "G", "", "Personal github token for administrator logins")
}

func main() {
	cmdMain.Execute()
}

func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func Exitf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf(format, args...)
	fmt.Println()
	os.Exit(1)
}

func Verbosef(format string, args ...interface{}) {
	if globalFlags.verbose {
		fmt.Printf(format, args...)
	}
}

func assert(err error) {
	if err != nil {
		Exitf("Assertion failed: %#v", err)
	}
}

func assertArgIsSet(arg, argKey string) {
	if arg == "" {
		Exitf("%s must be set\n", argKey)
	}
}
