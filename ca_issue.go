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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pulcy/vault-monkey/service"
	"github.com/spf13/cobra"
)

var (
	cmdCAIssue = &cobra.Command{
		Use:   "issue",
		Short: "Create vault information for a new cluster specific CA",
		Run:   showUsage,
	}

	cmdCAIssueETCD = &cobra.Command{
		Use:   "etcd",
		Short: "Issue a certificate to access the ETCD cluster.",
		Run:   cmdCAIssueETCDRun,
	}

	cmdCAIssueK8s = &cobra.Command{
		Use:   "k8s",
		Short: "Issue a certificate to access the Kubernetes API server.",
		Run:   cmdCAIssueK8sRun,
	}

	caIssueFlags struct {
		serverLogin bool
		service.IssueConfig
	}
)

func init() {
	cmdCA.AddCommand(cmdCAIssue)
	cmdCAIssue.AddCommand(cmdCAIssueETCD)
	cmdCAIssue.AddCommand(cmdCAIssueK8s)

	defaultOutputDir, err := homedir.Expand("~/.pulcy/certs")
	if err != nil {
		Exitf("Failed to expand homedir: %#v\n", err)
	}
	cmdCAIssue.PersistentFlags().BoolVar(&caIssueFlags.serverLogin, "server", false, "If set, a server login is performed (instead of admin login)")
	cmdCAIssue.PersistentFlags().StringVar(&caIssueFlags.IssueConfig.OutputDir, "destination", defaultOutputDir, "Where to store the issued certificates")
	cmdCAIssue.PersistentFlags().StringVar(&caIssueFlags.CommonName, "common-name", caIssueFlags.CommonName, "CommonName of the user")
	cmdCAIssue.PersistentFlags().StringVar(&caIssueFlags.Role, "role", caIssueFlags.Role, "Role used to issue certificate")
	cmdCAIssue.PersistentFlags().StringSliceVar(&caIssueFlags.AltNames, "alt-name", caIssueFlags.AltNames, "Alternate names of the user")
	cmdCAIssue.PersistentFlags().StringSliceVar(&caIssueFlags.IPSans, "ip-san", caIssueFlags.IPSans, "IP Subject Alternative Names")
	cmdCAIssue.PersistentFlags().StringVar(&caIssueFlags.CertificateFileName, "cert-file-name", caIssueFlags.CertificateFileName, "Filename of the public key")
	cmdCAIssue.PersistentFlags().StringVar(&caIssueFlags.KeyFileName, "key-file-name", caIssueFlags.KeyFileName, "Filename of the private key")
	cmdCAIssue.PersistentFlags().StringVar(&caIssueFlags.CAFileName, "ca-file-name", caIssueFlags.CAFileName, "Filename of the CA certificate")
	cmdCAIssue.PersistentFlags().Uint32Var(&caIssueFlags.FileMode, "file-mode", caIssueFlags.FileMode, "Mode of files that are created (defaults to 0600)")
}

func cmdCAIssueETCDRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")
	assertArgIsSet(caIssueFlags.OutputDir, "destination")
	assertArgIsSet(caIssueFlags.CommonName, "common-name")

	var c *service.AuthenticatedVaultClient
	var err error
	if caIssueFlags.serverLogin {
		c, _, err = serverLogin()
		if err != nil {
			Exitf("Server login failed: %v", err)
		}
	} else {
		_, c, err = adminLogin()
		if err != nil {
			Exitf("Admin login failed: %v", err)
		}
	}

	ca := c.CA()
	if err := ca.IssueETCDCertificate(caFlags.clusterID, caIssueFlags.IssueConfig); err != nil {
		Exitf("Failed to issue certificate: %v", err)
	}
}

func cmdCAIssueK8sRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")
	assertArgIsSet(caIssueFlags.OutputDir, "destination")
	assertArgIsSet(caIssueFlags.CommonName, "common-name")

	var c *service.AuthenticatedVaultClient
	var err error
	if caIssueFlags.serverLogin {
		c, _, err = serverLogin()
		if err != nil {
			Exitf("Server login failed: %v", err)
		}
	} else {
		_, c, err = adminLogin()
		if err != nil {
			Exitf("Admin login failed: %v", err)
		}
	}

	ca := c.CA()
	if err := ca.IssueK8sCertificate(caFlags.clusterID, caIssueFlags.IssueConfig); err != nil {
		Exitf("Failed to issue certificate: %v", err)
	}
}
