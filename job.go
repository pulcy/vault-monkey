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

	cmdJobAllowCluster = &cobra.Command{
		Use:   "allow",
		Short: "Allow a cluster to access secrets for a job",
		Run:   cmdJobAllowClusterRun,
	}

	cmdJobDenyCluster = &cobra.Command{
		Use:   "deny",
		Short: "Deny a cluster to access secrets for a job",
		Run:   cmdJobDenyClusterRun,
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
	cmdJob.AddCommand(cmdJobAllowCluster)
	cmdJob.AddCommand(cmdJobDenyCluster)

	cmdJob.PersistentFlags().StringVarP(&jobFlags.jobID, "job-id", "j", "", "ID of the job")
	cmdJob.PersistentFlags().StringVarP(&jobFlags.clusterID, "cluster-id", "c", "", "ID of the cluster")
	cmdJob.PersistentFlags().StringVarP(&jobFlags.policyName, "policy", "p", "", "Name of the policy for the job")
	cmdMain.AddCommand(cmdJob)
}

func cmdJobCreateRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")
	assertArgIsSet(jobFlags.policyName, "policy")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	job := c.Job()
	if err := job.Create(jobFlags.jobID, jobFlags.policyName); err != nil {
		Exitf("Failed to create job: %v", err)
	}
}

func cmdJobDeleteRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	job := c.Job()
	if err := job.Delete(jobFlags.jobID); err != nil {
		Exitf("Failed to delete job: %v", err)
	}
}

func cmdJobAllowClusterRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")
	assertArgIsSet(jobFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	job := c.Job()
	if err := job.AllowCluster(jobFlags.jobID, jobFlags.clusterID); err != nil {
		Exitf("Failed to allow cluster to access secrets of a job: %v", err)
	}
}

func cmdJobDenyClusterRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(jobFlags.jobID, "job-id")
	assertArgIsSet(jobFlags.clusterID, "cluster-id")

	_, c, err := adminLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	job := c.Job()
	if err := job.DenyCluster(jobFlags.jobID, jobFlags.clusterID); err != nil {
		Exitf("Failed to deny cluster access to secrets of a job: %v", err)
	}
}
