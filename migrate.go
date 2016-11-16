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
	"github.com/pulcy/vault-monkey/service"
	"github.com/pulcy/vault-monkey/service/migration"
	"github.com/spf13/cobra"
)

var (
	fromType    string
	fromAddress string
	toType      string
	toAddress   string
	cmdMigrate  = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate vault data between backends.",
		Run:   cmdMigrateRun,
	}
)

func init() {
	cmdMigrate.Flags().StringVar(&fromType, "from-type", "", "Type of source backend")
	cmdMigrate.Flags().StringVar(&fromAddress, "from-address", "", "Address of source backend")
	cmdMigrate.Flags().StringVar(&toType, "to-type", "", "Type of destination backend")
	cmdMigrate.Flags().StringVar(&toAddress, "to-address", "", "Address of destination backend")

	cmdMain.AddCommand(cmdMigrate)
}

func cmdMigrateRun(cmd *cobra.Command, args []string) {
	from := mustCreateBackend(fromType, fromAddress)
	to := mustCreateBackend(toType, toAddress)
	if err := service.Migrate(from, to, log); err != nil {
		Exitf("Failed to migrate vault data: %#v", err)
	}
}

func mustCreateBackend(bType, bAddress string) migration.Backend {
	var b migration.Backend
	var err error
	switch bType {
	case "etcd":
		b, err = migration.NewEtcdBackend(bAddress)
	default:
		Exitf("Unknown backend type '%s'", bType)
	}
	if err != nil {
		Exitf("Failed to create %s backend: %#v", bType, err)
	}
	return b
}
