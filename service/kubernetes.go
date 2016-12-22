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

package service

import (
	"fmt"
	"io/ioutil"
	"strings"

	k8s "github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
)

type K8sClient struct {
	baseServerLoginData
	c                     k8s.Client
	namespace             string
	podName               string
	podIP                 string
	clusterInfoSecretName string
	clusterIDSecretKey    string
}

// NewKubernetesClient creates a kubernetes client.
func NewKubernetesClient(podName, podIP, clusterInfoSecretName, clusterIDSecretKey string) (*K8sClient, error) {
	namespace, err := getKubernetesNamespace()
	if err != nil {
		return nil, maskAny(err)
	}

	client, err := http.NewInCluster()
	if err != nil {
		return nil, maskAny(err)
	}
	return &K8sClient{
		c:         client,
		namespace: namespace,
		podName:   podName,
		podIP:     podIP,
		clusterInfoSecretName: clusterInfoSecretName,
		clusterIDSecretKey:    clusterIDSecretKey,
	}, nil
}

func (c *K8sClient) ServerLoginData(next ServerLoginData) ServerLoginData {
	c.baseServerLoginData.next = next
	return c
}

func (c *K8sClient) JobID() (string, error) {
	// JobID is never read from Kubernetes secret
	return c.baseServerLoginData.JobID()
}

func (c *K8sClient) ClusterID() (string, error) {
	if c.clusterInfoSecretName != "" && c.clusterIDSecretKey != "" {
		s, err := c.getKubernetesSecret(c.clusterInfoSecretName)
		if err != nil {
			fmt.Printf("ClusterInfo secret with name '%s' not found", c.clusterInfoSecretName)
			// Now fallback to next
		} else {
			v, found := s.Data[c.clusterIDSecretKey]
			if !found {
				return "", maskAny(fmt.Errorf("Key '%s' is not found in secret '%s'", c.clusterIDSecretKey, c.clusterInfoSecretName))
			}
			return string(v), nil
		}
	}
	return c.baseServerLoginData.ClusterID()
}

func (c *K8sClient) MachineID() (string, error) {
	if c.podName != "" || c.podIP != "" {
		var hostIP string
		if c.podName != "" {
			if pod, err := c.c.GetPod(c.namespace, c.podName); err == nil {
				hostIP = pod.Status.HostIP
			}
		}
		if hostIP == "" {
			// This is the case then hostNetwork=true
			hostIP = c.podIP
		}
		if hostIP != "" {
			nodes, err := c.c.ListNodes(nil)
			if err != nil {
				return "", maskAny(err)
			}
			// Check by address
			for _, n := range nodes.Items {
				for _, a := range n.Status.Addresses {
					fmt.Printf("Checking node '%s' with address '%s', searching for '%s'\n", n.Spec.ExternalID, a.Address, hostIP)
					if a.Address == hostIP {
						nodeInfo := n.Status.NodeInfo
						id := nodeInfo.MachineID
						if id == "" {
							id = nodeInfo.SystemUUID
						}
						return id, nil
					}
				}
			}
			// Check by node name
			if c.podName != "" {
				for _, n := range nodes.Items {
					if n.Spec.ExternalID == c.podName {
						nodeInfo := n.Status.NodeInfo
						id := nodeInfo.MachineID
						if id == "" {
							id = nodeInfo.SystemUUID
						}
						return id, nil
					}
				}
			}
		}
	}
	return c.baseServerLoginData.MachineID()
}

// getKubernetesNamespace reads the namespace of the current pod from the well known location.
func getKubernetesNamespace() (string, error) {
	raw, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", maskAny(err)
	}
	return strings.TrimSpace(string(raw)), nil
}

func (c *K8sClient) getKubernetesSecret(secretName string) (k8s.Secret, error) {
	s, err := c.c.GetSecret(c.namespace, secretName)
	if err != nil {
		return k8s.Secret{}, maskAny(err)
	}
	return *s, nil
}

func (c *K8sClient) setKubernetesSecret(secretName string, secret k8s.Secret, create bool) error {
	if create {
		if _, err := c.c.CreateSecret(c.namespace, &secret); err != nil {
			return maskAny(err)
		}
	} else {
		if _, err := c.c.UpdateSecret(c.namespace, &secret); err != nil {
			return maskAny(err)
		}
	}

	return nil
}
