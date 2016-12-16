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
	"os"

	"github.com/spf13/cobra"

	"github.com/pulcy/vault-monkey/service"
)

const (
	defaultClusterIDPath            = "/etc/pulcy/cluster-id"
	defaultMachineIDPath            = "/etc/machine-id"
	defaultK8sClusterInfoSecretName = "vault-monkey-cluster-info"
	defaultK8sClusterIDSecretKey    = "CLUSTER_ID"
)

var (
	cmdExtract = &cobra.Command{
		Use:   "extract",
		Short: "Server commands to extract secrects from the vault",
		Run:   showUsage,
	}

	extractFlags struct {
		targetFilePath           string
		k8sPodName               string
		k8sClusterInfoSecretName string
		k8sClusterIDSecretKey    string
		k8sSecretName            string
		k8sSecretKey             string
		jobID                    string
		clusterIDPath            string
		machineIDPath            string
	}
)

func init() {
	hostName := os.Getenv("HOSTNAME")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.targetFilePath, "target", "", "Path of target file")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.k8sPodName, "kubernetes-pod-name", hostName, "Name of Kubernetes pod this process is running in")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.k8sClusterInfoSecretName, "kubernetes-cluster-info-secret-name", defaultK8sClusterInfoSecretName, "Name of Kubernetes secret that holds the cluster ID")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.k8sClusterIDSecretKey, "kubernetes-cluster-id-secret-key", defaultK8sClusterIDSecretKey, "Key for the cluster ID secret identified by `kubernetes-cluster-info-secret-name`")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.k8sSecretName, "kubernetes-secret-name", "", "Name of Kubernetes secret to store extracted data into")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.k8sSecretKey, "kubernetes-secret-key", "", "Key inside Kubernetes secret to store extracted data into")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.jobID, "job-id", "", "Identifier for the current job")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.clusterIDPath, "cluster-id-path", defaultClusterIDPath, "Path of cluster-id file")
	cmdExtract.PersistentFlags().StringVar(&extractFlags.machineIDPath, "machine-id-path", defaultMachineIDPath, "Path of machine-id file")
	cmdMain.AddCommand(cmdExtract)
}

// serverLogin initialized a VaultServices and tries to perform a server login.
func serverLogin() (*service.AuthenticatedVaultClient, *service.K8sClient, error) {
	// Create service
	vs, err := service.NewVaultService(log, globalFlags.VaultServiceConfig)
	if err != nil {
		return nil, nil, maskAny(err)
	}

	var k8sclient *service.K8sClient
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if host != "" && port != "" && (extractFlags.k8sPodName != "" || extractFlags.k8sClusterInfoSecretName != "") {
		k8sclient, err = service.NewKubernetesClient(extractFlags.k8sPodName, extractFlags.k8sClusterInfoSecretName, extractFlags.k8sClusterIDSecretKey)
		if err != nil {
			return nil, nil, maskAny(err)
		}
	}

	// Perform server login
	fsData := service.NewFileSystemServerLoginData("", extractFlags.clusterIDPath, extractFlags.machineIDPath, nil)
	staticData := service.NewStaticServerLoginData(extractFlags.jobID, "", "", fsData)
	envData := service.NewEnvServerLoginData(staticData)
	loginData := envData
	if k8sclient != nil {
		loginData = k8sclient.ServerLoginData()
	}

	c, err := vs.ServerLogin(loginData)
	if err != nil {
		return nil, nil, maskAny(err)
	}
	return c, k8sclient, nil
}
