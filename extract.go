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

	"github.com/spf13/cobra"

	"github.com/pulcy/vault-monkey/service"
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
func serverLogin() (*service.AuthenticatedVaultClient, error) {
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
	vs, err := service.NewVaultService(log, globalFlags.VaultServiceConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	// Perform server login
	c, err := vs.ServerLogin(extractFlags.ServerLoginData)
	if err != nil {
		return nil, maskAny(err)
	}
	return c, nil
}
