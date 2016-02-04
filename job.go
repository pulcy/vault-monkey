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
	cmdJob = &cobra.Command{
		Use:   "job",
		Short: "Administrator commands to manipulate vault information for jobs",
		Run:   showUsage,
	}

	cmdJobCreate = &cobra.Command{
		Use:   "create",
		Short: "Create vault information for a new job",
		Run:   cmdJobCreateRun,
	}

	cmdJobDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete vault information for a job",
		Run:   cmdJobDeleteRun,
	}

	cmdJobAddCluster = &cobra.Command{
		Use:   "add",
		Short: "Add a cluster to a job",
		Run:   cmdJobAddClusterRun,
	}

	cmdJobRemoveCluster = &cobra.Command{
		Use:   "remove",
		Short: "Remove a cluster from a job",
		Run:   cmdJobRemoveClusterRun,
	}

	jobFlags struct {
		jobID      string
		clusterID  string
		policyName string
	}
)

func init() {
	cmdJob.AddCommand(cmdJobCreate)
	cmdJob.AddCommand(cmdJobDelete)
	cmdJob.AddCommand(cmdJobAddCluster)
	cmdJob.AddCommand(cmdJobRemoveCluster)

	cmdJob.PersistentFlags().StringVarP(&jobFlags.jobID, "job-id", "j", "", "ID of the job")
	cmdJob.PersistentFlags().StringVarP(&jobFlags.clusterID, "cluster-id", "c", "", "ID of the cluster")
	cmdJob.PersistentFlags().StringVarP(&jobFlags.policyName, "policy", "p", "", "Name of the policy for the job")
	cmdMain.AddCommand(cmdJob)
}

func cmdJobCreateRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")
	assertArgIsSet(jobFlags.policyName, "policy")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Job().Create(jobFlags.jobID, jobFlags.policyName); err != nil {
		Exitf("Failed to create job: %v", err)
	}
}

func cmdJobDeleteRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Job().Delete(jobFlags.jobID); err != nil {
		Exitf("Failed to delete job: %v", err)
	}
}

func cmdJobAddClusterRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")
	assertArgIsSet(jobFlags.clusterID, "cluster-id")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Job().AddCluster(jobFlags.jobID, jobFlags.clusterID); err != nil {
		Exitf("Failed to add cluster to job: %v", err)
	}
}

func cmdJobRemoveClusterRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")
	assertArgIsSet(jobFlags.clusterID, "cluster-id")

	vs, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	if err := vs.Job().RemoveCluster(jobFlags.jobID, jobFlags.clusterID); err != nil {
		Exitf("Failed to remove cluster from job: %v", err)
	}
}
