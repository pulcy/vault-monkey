package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"./service"
)

const (
	defaultClusterIDPath = "/etc/pulcy/cluster-id"
	defaultMachineIDPath = "/etc/machine-id"
)

var (
	cmdExtract = &cobra.Command{
		Use:   "extract",
		Short: "Server commands to extract secrects from the vault",
		Run:   showUsage,
	}

	extractFlags struct {
		targetFilePath string
		service.ServerLoginData
	}
)

func init() {
	cmdExtract.PersistentFlags().StringVar(&extractFlags.targetFilePath, "target", "", "Path of target file")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.JobID, "job-id", "", "Identifier for the current job")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.ClusterIDPath, "cluster-id-path", defaultClusterIDPath, "Path of cluster-id file")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.MachineIDPath, "machine-id-path", defaultMachineIDPath, "Path of machine-id file")
	cmdMain.AddCommand(cmdExtract)
}

// serverLogin initialized a VaultServices and tries to perform a server login.
func serverLogin() (*service.VaultService, error) {
	// Check arguments
	if extractFlags.JobID == "" {
		return nil, maskAny(fmt.Errorf("--job-id missing"))
	}
	if extractFlags.ClusterIDPath == "" {
		return nil, maskAny(fmt.Errorf("--cluster-id-path missing"))
	}
	if extractFlags.MachineIDPath == "" {
		return nil, maskAny(fmt.Errorf("--machine-id-path missing"))
	}

	// Create service
	vs, err := service.NewVaultService(globalFlags.VaultServiceConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	// Perform server login
	if err := vs.ServerLogin(extractFlags.ServerLoginData); err != nil {
		return nil, maskAny(err)
	}
	return vs, nil
}
