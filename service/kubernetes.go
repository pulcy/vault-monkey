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
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/api/v1"
)

type K8sClient struct {
	baseServerLoginData
	c                     *k8s.Client
	namespace             string
	podName               string
	clusterInfoSecretName string
	clusterIDSecretKey    string
}

// NewKubernetesClient creates a kubernetes client.
func NewKubernetesClient(podName, clusterInfoSecretName, clusterIDSecretKey string) (*K8sClient, error) {
	namespace, err := getKubernetesNamespace()
	if err != nil {
		return nil, maskAny(err)
	}

	client, err := k8s.InClusterClient()
	if err != nil {
		return nil, maskAny(err)
	}
	return &K8sClient{
		c:                     client,
		namespace:             namespace,
		podName:               podName,
		clusterInfoSecretName: clusterInfoSecretName,
		clusterIDSecretKey:    clusterIDSecretKey,
	}, nil
}

func (c *K8sClient) ServerLoginData() ServerLoginData {
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
			v, found := s.GetData()[c.clusterIDSecretKey]
			if !found {
				return "", maskAny(fmt.Errorf("Key '%s' is not found in secret '%s'", c.clusterIDSecretKey, c.clusterInfoSecretName))
			}
			return string(v), nil
		}
	}
	return c.baseServerLoginData.ClusterID()
}

func (c *K8sClient) MachineID() (string, error) {
	if c.podName != "" {
		ctx := k8s.NamespaceContext(context.Background(), c.namespace)
		pod, err := c.c.CoreV1().GetPod(ctx, c.podName)
		if err != nil {
			return "", maskAny(err)
		}
		podHostIP := pod.GetStatus().GetHostIP()
		nodes, err := c.c.CoreV1().ListNodes(ctx)
		if err != nil {
			return "", maskAny(err)
		}
		for _, n := range nodes.Items {
			for _, a := range n.GetStatus().GetAddresses() {
				if a.GetAddress() == podHostIP {
					return n.GetStatus().GetNodeInfo().GetMachineID(), nil
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

func (c *K8sClient) getKubernetesSecret(secretName string) (v1.Secret, error) {
	ctx := k8s.NamespaceContext(context.Background(), c.namespace)
	s, err := c.c.CoreV1().GetSecret(ctx, secretName)
	if err != nil {
		return v1.Secret{}, maskAny(err)
	}
	return *s, nil
}

func (c *K8sClient) setKubernetesSecret(secretName string, secret v1.Secret, create bool) error {
	ctx := k8s.NamespaceContext(context.Background(), c.namespace)
	api := c.c.CoreV1()
	if create {
		if _, err := api.CreateSecret(ctx, &secret); err != nil {
			return maskAny(err)
		}
	} else {
		if _, err := api.UpdateSecret(ctx, &secret); err != nil {
			return maskAny(err)
		}
	}

	return nil
}
