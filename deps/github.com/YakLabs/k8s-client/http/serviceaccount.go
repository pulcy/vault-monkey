package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func serviceaccountGeneratePath(namespace, name string) string {
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/serviceaccounts"
	}
	return "/api/v1/namespaces/" + namespace + "/serviceaccounts/" + name
}

// GetServiceAccount fetches a single ServiceAccount
func (c *Client) GetServiceAccount(namespace, name string) (*k8s.ServiceAccount, error) {
	var out k8s.ServiceAccount
	_, err := c.do("GET", serviceaccountGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ServiceAccount")
	}
	return &out, nil
}

// CreateServiceAccount creates a new ServiceAccount. This will fail if it already exists.
func (c *Client) CreateServiceAccount(namespace string, item *k8s.ServiceAccount) (*k8s.ServiceAccount, error) {
	item.TypeMeta.Kind = "ServiceAccount"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ServiceAccount
	_, err := c.do("POST", serviceaccountGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ServiceAccount")
	}
	return &out, nil
}

// ListServiceAccounts lists all ServiceAccounts in a namespace
func (c *Client) ListServiceAccounts(namespace string, opts *k8s.ListOptions) (*k8s.ServiceAccountList, error) {
	var out k8s.ServiceAccountList
	_, err := c.do("GET", serviceaccountGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ServiceAccounts")
	}
	return &out, nil
}

// DeleteServiceAccount deletes a single ServiceAccount. It will error if the ServiceAccount does not exist.
func (c *Client) DeleteServiceAccount(namespace, name string) error {
	_, err := c.do("DELETE", serviceaccountGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete ServiceAccount")
}

// UpdateServiceAccount will update in place a single ServiceAccount. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateServiceAccount(namespace string, item *k8s.ServiceAccount) (*k8s.ServiceAccount, error) {
	item.TypeMeta.Kind = "ServiceAccount"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ServiceAccount
	_, err := c.do("PUT", serviceaccountGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update ServiceAccount")
	}
	return &out, nil
}
