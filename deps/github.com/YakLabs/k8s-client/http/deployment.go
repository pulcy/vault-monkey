package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func deploymentGeneratePath(namespace, name string) string {
	if name == "" {
		return "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments"
	}
	return "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments/" + name
}

// GetDeployment fetches a single Deployment
func (c *Client) GetDeployment(namespace, name string) (*k8s.Deployment, error) {
	var out k8s.Deployment
	_, err := c.do("GET", deploymentGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Deployment")
	}
	return &out, nil
}

// CreateDeployment creates a new Deployment. This will fail if it already exists.
func (c *Client) CreateDeployment(namespace string, item *k8s.Deployment) (*k8s.Deployment, error) {
	item.TypeMeta.Kind = "Deployment"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Deployment
	_, err := c.do("POST", deploymentGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Deployment")
	}
	return &out, nil
}

// ListDeployments lists all Deployments in a namespace
func (c *Client) ListDeployments(namespace string, opts *k8s.ListOptions) (*k8s.DeploymentList, error) {
	var out k8s.DeploymentList
	_, err := c.do("GET", deploymentGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Deployments")
	}
	return &out, nil
}

// DeleteDeployment deletes a single Deployment. It will error if the Deployment does not exist.
func (c *Client) DeleteDeployment(namespace, name string) error {
	_, err := c.do("DELETE", deploymentGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Deployment")
}

// UpdateDeployment will update in place a single Deployment. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateDeployment(namespace string, item *k8s.Deployment) (*k8s.Deployment, error) {
	item.TypeMeta.Kind = "Deployment"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Deployment
	_, err := c.do("PUT", deploymentGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Deployment")
	}
	return &out, nil
}
