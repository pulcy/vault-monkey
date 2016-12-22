package client

type (
	ConfigMapInterface interface {
		CreateConfigMap(namespace string, item *ConfigMap) (*ConfigMap, error)
		GetConfigMap(namespace, name string) (result *ConfigMap, err error)
		ListConfigMaps(namespace string, opts *ListOptions) (*ConfigMapList, error)
		DeleteConfigMap(namespace, name string) error
		UpdateConfigMap(namespace string, item *ConfigMap) (*ConfigMap, error)
	}

	ConfigMapType string

	ConfigMap struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Data       map[string][]byte `json:"data,omitempty"`
	}

	ConfigMapList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []ConfigMap `json:"items"`
	}
)
