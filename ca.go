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
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cmdCA = &cobra.Command{
		Use:              "ca",
		Short:            "Administrator commands to create CA's and issue certificates",
		Run:              showUsage,
		PersistentPreRun: cmdCAPersistentPreRun,
	}

	cmdCACreate = &cobra.Command{
		Use:   "create",
		Short: "Create a cluster specific CA",
		Run:   showUsage,
	}

	cmdCACreateETCD = &cobra.Command{
		Use:   "etcd",
		Short: "Create vault CA that is to be used by ETCD members",
		Run:   cmdCACreateETCDRun,
	}

	cmdCACreateK8s = &cobra.Command{
		Use:   "k8s",
		Short: "Create vault CA that is to be used by Kubernetes components",
		Run:   cmdCACreateK8sRun,
	}

	caFlags struct {
		clusterID     string
		clusterIDFile string
		force         bool
		component     string
		domainName    string
	}
)

func init() {
	cmdMain.AddCommand(cmdCA)
	cmdCA.AddCommand(cmdCACreate)
	cmdCACreate.AddCommand(cmdCACreateETCD)
	cmdCACreate.AddCommand(cmdCACreateK8s)

	cmdCA.PersistentFlags().StringVar(&caFlags.clusterID, "cluster-id", "", "ID of the cluster to create a CA for")
	cmdCA.PersistentFlags().StringVar(&caFlags.clusterIDFile, "cluster-id-file", "", "Path of the file containing cluster-id")
	cmdCACreate.PersistentFlags().BoolVar(&caFlags.force, "force", false, "If set, existing mounts will be overwritten, revoking issues certificates")
	cmdCACreate.PersistentFlags().StringVar(&caFlags.domainName, "domain", "", "Domain name of the cluster")
	cmdCACreateK8s.Flags().StringVar(&caFlags.component, "component", "", "The Kubernetes component name to create a CA for")
}

func cmdCAPersistentPreRun(cmd *cobra.Command, args []string) {
	if caFlags.clusterID == "" && caFlags.clusterIDFile != "" {
		raw, err := ioutil.ReadFile(caFlags.clusterIDFile)
		if err != nil {
			Exitf("Failed to read %s: %#v\n", caFlags.clusterIDFile, err)
		}
		caFlags.clusterID = strings.TrimSpace(string(raw))
	}
}

func cmdCACreateETCDRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")
	assertArgIsSet(caFlags.domainName, "domain")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.CreateETCDMembers(caFlags.clusterID, caFlags.domainName, caFlags.force); err != nil {
		Exitf("Failed to create CA: %v", err)
	}
}

func cmdCACreateK8sRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")
	assertArgIsSet(caFlags.domainName, "domain")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if caFlags.component != "" {
		if err := ca.CreateK8s(caFlags.clusterID, caFlags.component, caFlags.domainName, caFlags.force); err != nil {
			Exitf("Failed to create CA: %v", err)
		}
	} else {
		if err := ca.CreateK8sAll(caFlags.clusterID, caFlags.domainName, caFlags.force); err != nil {
			Exitf("Failed to create CA's: %v", err)
		}
	}
}
