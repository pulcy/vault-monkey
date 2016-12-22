package client

type (
	// ServiceAccountInterface is the interface that defines ServiceAccount functions.
	ServiceAccountInterface interface {
		CreateServiceAccount(namespace string, item *ServiceAccount) (*ServiceAccount, error)
		GetServiceAccount(namespace, name string) (result *ServiceAccount, err error)
		ListServiceAccounts(namespace string, opts *ListOptions) (*ServiceAccountList, error)
		DeleteServiceAccount(namepsace, name string) error
		UpdateServiceAccount(namespace string, item *ServiceAccount) (*ServiceAccount, error)
	}

	// ServiceAccount binds together: * a name, understood by users, and perhaps by peripheral systems, for an identity * a principal that can be authenticated and authorized * a set of secrets
	ServiceAccount struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// Secrets is the list of secrets allowed to be used by pods running using this ServiceAccount
		Secrets []ObjectReference `json:"secrets"`

		// ImagePullSecrets is a list of references to secrets in the same namespace to use for pulling any images
		// in pods that reference this ServiceAccount.  ImagePullSecrets are distinct from Secrets because Secrets
		// can be mounted in the pod, but ImagePullSecrets are only accessed by the kubelet.
		ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	}

	// ServiceAccountList is a list of ServiceAccount objects
	ServiceAccountList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		Items    []ServiceAccount `json:"items"`
	}
)
