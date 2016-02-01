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

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Cluster().Create(clusterFlags.clusterID); err != nil {
		Exitf("Failed to create cluster: %v", err)
	}
}

func cmdClusterDeleteRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Cluster().Delete(clusterFlags.clusterID); err != nil {
		Exitf("Failed to create cluster: %v", err)
	}
}

func cmdClusterAddMachineRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.clusterID, "cluster-id")
	assertArgIsSet(clusterFlags.machineID, "machine-id")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Cluster().AddMachine(clusterFlags.clusterID, clusterFlags.machineID, clusterFlags.cidrBlock); err != nil {
		Exitf("Failed to add machine to cluster: %v", err)
	}
}

func cmdClusterRemoveMachineRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(clusterFlags.machineID, "machine-id")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Cluster().RemoveMachine(clusterFlags.machineID); err != nil {
		Exitf("Failed to remove machine from cluster: %v", err)
	}
}
