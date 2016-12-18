package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func configmapGeneratePath(namespace, name string) string {
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/configmaps"
	}
	return "/api/v1/namespaces/" + namespace + "/configmaps/" + name
}

// GetConfigMap fetches a single ConfigMap
func (c *Client) GetConfigMap(namespace, name string) (*k8s.ConfigMap, error) {
	var out k8s.ConfigMap
	_, err := c.do("GET", configmapGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ConfigMap")
	}
	return &out, nil
}

// CreateConfigMap creates a new ConfigMap. This will fail if it already exists.
func (c *Client) CreateConfigMap(namespace string, item *k8s.ConfigMap) (*k8s.ConfigMap, error) {
	item.TypeMeta.Kind = "ConfigMap"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ConfigMap
	_, err := c.do("POST", configmapGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ConfigMap")
	}
	return &out, nil
}

// ListConfigMaps lists all ConfigMaps in a namespace
func (c *Client) ListConfigMaps(namespace string, opts *k8s.ListOptions) (*k8s.ConfigMapList, error) {
	var out k8s.ConfigMapList
	_, err := c.do("GET", configmapGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ConfigMaps")
	}
	return &out, nil
}

// DeleteConfigMap deletes a single ConfigMap. It will error if the ConfigMap does not exist.
func (c *Client) DeleteConfigMap(namespace, name string) error {
	_, err := c.do("DELETE", configmapGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete ConfigMap")
}

// UpdateConfigMap will update in place a single ConfigMap. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateConfigMap(namespace string, item *k8s.ConfigMap) (*k8s.ConfigMap, error) {
	item.TypeMeta.Kind = "ConfigMap"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ConfigMap
	_, err := c.do("PUT", configmapGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update ConfigMap")
	}
	return &out, nil
}
