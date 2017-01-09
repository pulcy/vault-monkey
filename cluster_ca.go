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
	"github.com/spf13/cobra"
)

var (
	cmdClusterCreateCA = &cobra.Command{
		Use:   "create-ca",
		Short: "Create vault information for a new cluster specific CA",
		Run:   showUsage,
	}

	cmdClusterCreateCAETCD = &cobra.Command{
		Use:   "etcd",
		Short: "Create vault CA that is to be used by ETCD members",
		Run:   cmdClusterCreateCAETCDRun,
	}

	cmdClusterCreateCAK8s = &cobra.Command{
		Use:   "k8s",
		Short: "Create vault CA that is to be used by Kubernetes components",
		Run:   cmdClusterCreateCAK8sRun,
	}

	caFlags struct {
		force     bool
		component string
	}
)

func init() {
	cmdCluster.AddCommand(cmdClusterCreateCA)
	cmdClusterCreateCA.AddCommand(cmdClusterCreateCAETCD)
	cmdClusterCreateCA.AddCommand(cmdClusterCreateCAK8s)

	cmdClusterCreateCA.PersistentFlags().BoolVar(&caFlags.force, "force", false, "If set, existing mounts will be overwritten, revoking issues certificates")
	cmdClusterCreateCAK8s.Flags().StringVar(&caFlags.component, "component", "", "The Kubernetes component name to create a CA for")
}

func cmdClusterCreateCAETCDRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.CreateETCDMembers(clusterFlags.clusterID, caFlags.force); err != nil {
		Exitf("Failed to create CA: %v", err)
	}
}

func cmdClusterCreateCAK8sRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if caFlags.component != "" {
		if err := ca.CreateK8s(clusterFlags.clusterID, caFlags.component, caFlags.force); err != nil {
			Exitf("Failed to create CA: %v", err)
		}
	} else {
		if err := ca.CreateK8sAll(clusterFlags.clusterID, caFlags.force); err != nil {
			Exitf("Failed to create CA's: %v", err)
		}
	}
}
