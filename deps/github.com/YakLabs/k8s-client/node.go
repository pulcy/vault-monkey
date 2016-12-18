package client

type (
	NodeInterface interface {
		CreateNode(item *Node) (*Node, error)
		GetNode(name string) (result *Node, err error)
		ListNodes(opts *ListOptions) (*NodeList, error)
		DeleteNode(name string) error
		UpdateNode(item *Node) (*Node, error)
	}

	NodeSpec struct {
		PodCIDR       string `json:"podCIDR,omitempty"`
		ExternalID    string `json:"externalID,omitempty"`
		ProviderID    string `json:"providerID,omitempty"`
		Unschedulable bool   `json:"unschedulable,omitempty"`
	}

	NodePhase         string
	NodeConditionType string
	NodeAddressType   string

	NodeAddress struct {
		Type    NodeAddressType `json:"type"`
		Address string          `json:"address"`
	}

	NodeCondition struct {
		Type    NodeConditionType `json:"type"`
		Status  ConditionStatus   `json:"status"`
		Reason  string            `json:"reason,omitempty"`
		Message string            `json:"message,omitempty"`
	}

	NodeStatus struct {
		Phase      NodePhase       `json:"phase,omitempty"`
		Conditions []NodeCondition `json:"conditions,omitempty"`
		Addresses  []NodeAddress   `json:"addresses,omitempty"`
		NodeInfo   *NodeSystemInfo `json:"nodeInfo,omitempty"`
	}

	// NodeSystemInfo is a set of ids/uuids to uniquely identify the node.
	NodeSystemInfo struct {
		// MachineID reported by the node. For unique machine identification in the cluster this field is prefered.
		MachineID string `json:"machineID"`
		// SystemUUID reported by the node. For unique machine identification MachineID is prefered. This field is specific to Red Hat hosts
		SystemUUID string `json:"systemUUID"`
		// Boot ID reported by the node.
		BootID string `json:"bootID"`
		// Kernel Version reported by the node from uname -r (e.g. 3.16.0-0.bpo.4-amd64).
		KernelVersion string `json:"kernelVersion"`
		// OS Image reported by the node from /etc/os-release (e.g. Debian GNU/Linux 7 (wheezy)).
		OSImage string `json:"osImage"`
		// ContainerRuntime Version reported by the node through runtime remote API (e.g. docker://1.5.0).
		ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
		// Kubelet Version reported by the node.
		KubeletVersion string `json:"kubeletVersion"`
		// KubeProxy Version reported by the node.
		KubeProxyVersion string `json:"kubeProxyVersion"`
		// The Operating System reported by the node.
		OperatingSystem string `json:"operatingSystem"`
		// The Architecture reported by the node.
		Architecture string `json:"architecture"`
	}

	Node struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Spec       NodeSpec   `json:"spec,omitempty"`
		Status     NodeStatus `json:"status,omitempty"`
	}

	NodeList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		Items    []Node `json:"items"`
	}
)
