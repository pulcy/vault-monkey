package client

type (
	NamespaceInterface interface {
		CreateNamespace(item *Namespace) (*Namespace, error)
		GetNamespace(name string) (result *Namespace, err error)
		ListNamespaces(opts *ListOptions) (*NamespaceList, error)
		DeleteNamespace(name string) error
		UpdateNamespace(item *Namespace) (*Namespace, error)
	}

	NamespaceSpec struct {
		Finalizers []FinalizerName
	}

	NamespacePhase string

	NamespaceStatus struct {
		Phase NamespacePhase `json:"phase,omitempty"`
	}

	Namespace struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Spec       NamespaceSpec   `json:"spec,omitempty"`
		Status     NamespaceStatus `json:"status,omitempty"`
	}

	NamespaceList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []Namespace `json:"items"`
	}
)

// NewNamespace creates a new namespace struct
func NewNamespace(name string) *Namespace {
	return &Namespace{
		TypeMeta: TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMeta{
			Namespace:   name,
			Name:        name,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
	}

}
