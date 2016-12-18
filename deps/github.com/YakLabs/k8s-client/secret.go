package client

type (
	// SecretInterface has methods to work with Secret resources.
	SecretInterface interface {
		CreateSecret(namespace string, item *Secret) (*Secret, error)
		GetSecret(namespace, name string) (result *Secret, err error)
		ListSecrets(namespace string, opts *ListOptions) (*SecretList, error)
		DeleteSecret(namespace, name string) error
		UpdateSecret(namespace string, item *Secret) (*Secret, error)
	}

	// SecretType is the type of secret.
	SecretType string

	// Secret holds secret data of a certain type.
	Secret struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Data       map[string][]byte `json:"data,omitempty"`

		// Used to facilitate programmatic handling of secret data.
		Type SecretType `json:"type,omitempty"`
	}

	// SecretList holds a list of secrets.
	SecretList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []Secret `json:"items"`
	}
)

// NewSecret creates a new Secret struct
func NewSecret(namespace, name string) *Secret {
	return &Secret{
		TypeMeta: TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Data: make(map[string][]byte),
	}
}
