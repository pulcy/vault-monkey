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
	"github.com/spf13/cobra"
)

var (
	cmdClusterIssueCertificate = &cobra.Command{
		Use:   "issue-cert",
		Short: "Create vault information for a new cluster specific CA",
		Run:   showUsage,
	}

	cmdClusterIssueK8sCertificate = &cobra.Command{
		Use:   "k8s",
		Short: "Issue an operation certificate to access the Kubernetes API server.",
		Run:   cmdClusterIssueK8sCertificateRun,
	}

	caIssueFlags struct {
		userName  string
		outputDir string
	}
)

func init() {
	cmdCluster.AddCommand(cmdClusterIssueCertificate)
	cmdClusterIssueCertificate.AddCommand(cmdClusterIssueK8sCertificate)

	defaultOutputDir, err := homedir.Expand("~/.pulcy/certs")
	if err != nil {
		Exitf("Failed to expand homedir: %#v\n", err)
	}
	cmdClusterIssueCertificate.PersistentFlags().StringVar(&caIssueFlags.outputDir, "destination", defaultOutputDir, "Where to store the issued certificates")
	cmdClusterIssueCertificate.PersistentFlags().StringVar(&caIssueFlags.userName, "username", "", "Name of the user")
}

func cmdClusterIssueK8sCertificateRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")
	assertArgIsSet(caIssueFlags.outputDir, "destination")
	assertArgIsSet(caIssueFlags.userName, "username")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.IssueK8sUserCertificate(clusterFlags.clusterID, caIssueFlags.userName, caIssueFlags.outputDir); err != nil {
		Exitf("Failed to issue certificate: %v", err)
	}
}
