package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func secretGeneratePath(namespace, name string) string {
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/secrets"
	}
	return "/api/v1/namespaces/" + namespace + "/secrets/" + name
}

// GetSecret fetches a single Secret
func (c *Client) GetSecret(namespace, name string) (*k8s.Secret, error) {
	var out k8s.Secret
	_, err := c.do("GET", secretGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Secret")
	}
	return &out, nil
}

// CreateSecret creates a new Secret. This will fail if it already exists.
func (c *Client) CreateSecret(namespace string, item *k8s.Secret) (*k8s.Secret, error) {
	item.TypeMeta.Kind = "Secret"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Secret
	_, err := c.do("POST", secretGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Secret")
	}
	return &out, nil
}

// ListSecrets lists all Secrets in a namespace
func (c *Client) ListSecrets(namespace string, opts *k8s.ListOptions) (*k8s.SecretList, error) {
	var out k8s.SecretList
	_, err := c.do("GET", secretGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Secrets")
	}
	return &out, nil
}

// DeleteSecret deletes a single Secret. It will error if the Secret does not exist.
func (c *Client) DeleteSecret(namespace, name string) error {
	_, err := c.do("DELETE", secretGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Secret")
}

// UpdateSecret will update in place a single Secret. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateSecret(namespace string, item *k8s.Secret) (*k8s.Secret, error) {
	item.TypeMeta.Kind = "Secret"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Secret
	_, err := c.do("PUT", secretGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Secret")
	}
	return &out, nil
}
