package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

// GetNamespace gets a namespace
func (c *Client) GetNamespace(name string) (*k8s.Namespace, error) {
	var out k8s.Namespace
	_, err := c.do("GET", "/api/v1/namespaces/"+name, nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get namespace")
	}
	return &out, nil
}

// CreateNamespace creates a new namespace. It will fail if the namespace already exists.
func (c *Client) CreateNamespace(item *k8s.Namespace) (*k8s.Namespace, error) {
	item.TypeMeta.Kind = "Namespace"
	item.TypeMeta.APIVersion = "v1"

	var out k8s.Namespace
	_, err := c.do("POST", "/api/v1/namespaces", item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create namespace")
	}
	return &out, nil
}

// ListNamespaces list all namespaces, optionally filtering.
func (c *Client) ListNamespaces(opts *k8s.ListOptions) (*k8s.NamespaceList, error) {
	var out k8s.NamespaceList
	_, err := c.do("GET", "/api/v1/namespaces?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list namespaces")
	}
	return &out, nil
}

// DeleteNamespace deletes a single namespace. It will error it it does not exist.
func (c *Client) DeleteNamespace(name string) error {
	_, err := c.do("DELETE", "/api/v1/namespaces/"+name, nil, nil)
	return errors.Wrap(err, "failed to delete namespace")
}

// UpdateNamespace updates a namespace.
func (c *Client) UpdateNamespace(item *k8s.Namespace) (*k8s.Namespace, error) {
	item.TypeMeta.Kind = "Namespace"
	item.TypeMeta.APIVersion = "v1"

	var out k8s.Namespace
	_, err := c.do("PUT", "/api/v1/namespaces/"+item.Name, item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update namespace")
	}
	return &out, nil
}
