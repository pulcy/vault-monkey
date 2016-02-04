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
	cmdSeal = &cobra.Command{
		Use:   "seal",
		Short: "Seal the vault.",
		Run:   cmdSealRun,
	}
)

func init() {
	cmdMain.AddCommand(cmdSeal)
}

func cmdSealRun(cmd *cobra.Command, args []string) {
	vs, err := adminLogin()
	if err != nil {
		Exitf("Failed to create vault service: %#v", err)
	}
	if err := vs.Seal(); err != nil {
		Exitf("Failed to seal vault: %#v", err)
	}
}
