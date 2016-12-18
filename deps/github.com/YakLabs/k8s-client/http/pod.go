package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func podGeneratePath(namespace, name string) string {
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/pods"
	}
	return "/api/v1/namespaces/" + namespace + "/pods/" + name
}

// GetPod fetches a single Pod
func (c *Client) GetPod(namespace, name string) (*k8s.Pod, error) {
	var out k8s.Pod
	_, err := c.do("GET", podGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Pod")
	}
	return &out, nil
}

// CreatePod creates a new Pod. This will fail if it already exists.
func (c *Client) CreatePod(namespace string, item *k8s.Pod) (*k8s.Pod, error) {
	item.TypeMeta.Kind = "Pod"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Pod
	_, err := c.do("POST", podGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Pod")
	}
	return &out, nil
}

// ListPods lists all Pods in a namespace
func (c *Client) ListPods(namespace string, opts *k8s.ListOptions) (*k8s.PodList, error) {
	var out k8s.PodList
	_, err := c.do("GET", podGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Pods")
	}
	return &out, nil
}

// DeletePod deletes a single Pod. It will error if the Pod does not exist.
func (c *Client) DeletePod(namespace, name string) error {
	_, err := c.do("DELETE", podGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Pod")
}

// UpdatePod will update in place a single Pod. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdatePod(namespace string, item *k8s.Pod) (*k8s.Pod, error) {
	item.TypeMeta.Kind = "Pod"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Pod
	_, err := c.do("PUT", podGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Pod")
	}
	return &out, nil
}
