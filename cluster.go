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
	"github.com/spf13/cobra"
)

var (
	cmdCluster = &cobra.Command{
		Use:   "cluster",
		Short: "Administator commands to manipulate vault information for clusters",
		Run:   showUsage,
	}

	cmdClusterCreate = &cobra.Command{
		Use:   "create",
		Short: "Create vault information for a new cluster",
		Run:   cmdClusterCreateRun,
	}

	cmdClusterDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete vault information for a cluster",
		Run:   cmdClusterDeleteRun,
	}

	cmdClusterAddMachine = &cobra.Command{
		Use:   "add",
		Short: "Add a machine to a cluster",
		Run:   cmdClusterAddMachineRun,
	}

	cmdClusterRemoveMachine = &cobra.Command{
		Use:   "remove",
		Short: "Remove a machine from a cluster",
		Run:   cmdClusterRemoveMachineRun,
	}

	clusterFlags struct {
		clusterID string
		machineID string
		cidrBlock string
	}
)

func init() {
	cmdCluster.AddCommand(cmdClusterCreate)
	cmdCluster.AddCommand(cmdClusterDelete)
	cmdCluster.AddCommand(cmdClusterAddMachine)
	cmdCluster.AddCommand(cmdClusterRemoveMachine)

	cmdCluster.PersistentFlags().StringVarP(&clusterFlags.clusterID, "cluster-id", "c", "", "ID of the cluster")
	cmdCluster.PersistentFlags().StringVarP(&clusterFlags.machineID, "machine-id", "m", "", "ID of the machine")
	cmdCluster.PersistentFlags().StringVar(&clusterFlags.cidrBlock, "cidr", "", "CIDR block from which the machine is allowed to connect")
	cmdMain.AddCommand(cmdCluster)
}

func cmdClusterCreateRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	cluster := c.Cluster()
	if err := cluster.Create(clusterFlags.clusterID); err != nil {
		Exitf("Failed to create cluster: %v", err)
	}
}

func cmdClusterDeleteRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	cluster := c.Cluster()
	if err := cluster.Delete(clusterFlags.clusterID); err != nil {
		Exitf("Failed to create cluster: %v", err)
	}
}

func cmdClusterAddMachineRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")
	assertArgIsSet(clusterFlags.machineID, "machine-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	cluster := c.Cluster()
	if err := cluster.AddMachine(clusterFlags.clusterID, clusterFlags.machineID, clusterFlags.cidrBlock); err != nil {
		Exitf("Failed to add machine to cluster: %v", err)
	}
}

func cmdClusterRemoveMachineRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.machineID, "machine-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	cluster := c.Cluster()
	if err := cluster.RemoveMachine(clusterFlags.machineID); err != nil {
		Exitf("Failed to remove machine from cluster: %v", err)
	}
}
