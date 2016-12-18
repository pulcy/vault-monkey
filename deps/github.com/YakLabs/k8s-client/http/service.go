package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func serviceGeneratePath(namespace, name string) string {
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/services"
	}
	return "/api/v1/namespaces/" + namespace + "/services/" + name
}

// GetService fetches a single Service
func (c *Client) GetService(namespace, name string) (*k8s.Service, error) {
	var out k8s.Service
	_, err := c.do("GET", serviceGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Service")
	}
	return &out, nil
}

// CreateService creates a new Service. This will fail if it already exists.
func (c *Client) CreateService(namespace string, item *k8s.Service) (*k8s.Service, error) {
	item.TypeMeta.Kind = "Service"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Service
	_, err := c.do("POST", serviceGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Service")
	}
	return &out, nil
}

// ListServices lists all Services in a namespace
func (c *Client) ListServices(namespace string, opts *k8s.ListOptions) (*k8s.ServiceList, error) {
	var out k8s.ServiceList
	_, err := c.do("GET", serviceGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Services")
	}
	return &out, nil
}

// DeleteService deletes a single Service. It will error if the Service does not exist.
func (c *Client) DeleteService(namespace, name string) error {
	_, err := c.do("DELETE", serviceGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Service")
}

// UpdateService will update in place a single Service. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateService(namespace string, item *k8s.Service) (*k8s.Service, error) {
	item.TypeMeta.Kind = "Service"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Service
	_, err := c.do("PUT", serviceGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Service")
	}
	return &out, nil
}
