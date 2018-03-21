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
	"time"

	"github.com/pulcy/vault-monkey/service"
	"github.com/spf13/cobra"
)

var (
	cmdCATidy = &cobra.Command{
		Use:   "tidy",
		Short: "Perform cleanup of expired certificates for a specific CA",
		Run:   showUsage,
	}

	cmdCATidyETCD = &cobra.Command{
		Use:   "etcd",
		Short: "Perform cleanup of expired certificates for a specific ETCD CA",
		Run:   cmdCATidyETCDRun,
	}

	cmdCATidyK8s = &cobra.Command{
		Use:   "k8s",
		Short: "Perform cleanup of expired certificates for a specific K8s CA",
		Run:   cmdCATidyK8sRun,
	}

	tidyFlags service.TidyOptions
)

func init() {
	f := cmdCATidy.PersistentFlags()
	f.BoolVar(&tidyFlags.TidyCertificateStore, "certificate-store", true, "If set, cleans the certificate store")
	f.BoolVar(&tidyFlags.TidyRevocationList, "revocation-list", true, "If set, cleans the certificate revocation list")
	f.DurationVar(&tidyFlags.SafetyBuffer, "safety-buffer", time.Hour*72, "Specifies a safety buffer to ensure certificates are not expunged prematurely")

	cmdCA.AddCommand(cmdCATidy)
	cmdCATidy.AddCommand(cmdCATidyETCD)
	cmdCATidy.AddCommand(cmdCATidyK8s)
}

func cmdCATidyETCDRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.TidyETCDCertificates(caFlags.clusterID, tidyFlags); err != nil {
		Exitf("Failed to list certificates: %v", err)
	}
}

func cmdCATidyK8sRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.TidyK8sCertificates(caFlags.clusterID, tidyFlags); err != nil {
		Exitf("Failed to list certificates: %v", err)
	}
}
