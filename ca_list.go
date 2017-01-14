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

import "github.com/spf13/cobra"

var (
	cmdCAList = &cobra.Command{
		Use:   "list",
		Short: "List vault information for a new cluster specific CA",
		Run:   showUsage,
	}

	cmdCAListK8s = &cobra.Command{
		Use:   "k8s",
		Short: "List certificates issued to access the Kubernetes API server.",
		Run:   cmdCAListK8sRun,
	}
)

func init() {
	cmdCA.AddCommand(cmdCAList)
	cmdCAList.AddCommand(cmdCAListK8s)
}

func cmdCAListK8sRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(caFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	ca := c.CA()
	if err := ca.ListK8sCertificates(caFlags.clusterID); err != nil {
		Exitf("Failed to list certificates: %v", err)
	}
}
