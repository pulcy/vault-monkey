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
	cmdCA = &cobra.Command{
		Use:   "ca",
		Short: "Administator commands to manipulate vault information for CA's",
		Run:   showUsage,
	}

	cmdCACreate = &cobra.Command{
		Use:   "create",
		Short: "Create vault information for a new CA",
		Run:   showUsage,
	}

	cmdCACreateETCD = &cobra.Command{
		Use:   "etcd",
		Short: "Create vault CA that is to be used by ETCD",
		Run:   cmdCACreateETCDRun,
	}

	caFlags struct {
		mountPoint string
	}
)

func init() {
	cmdCA.AddCommand(cmdCACreate)
	cmdCACreate.AddCommand(cmdCACreateETCD)

	cmdCACreate.PersistentFlags().StringVarP(&caFlags.mountPoint, "mountpoint", "m", "", "Mountpoint of the CA")

	cmdMain.AddCommand(cmdCA)
}

func cmdCACreateETCDRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.mountPoint, "mountpoint")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.CreateETCD(caFlags.mountPoint); err != nil {
		Exitf("Failed to create CA: %v", err)
	}
}
