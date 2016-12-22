package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func replicasetGeneratePath(namespace, name string) string {
	if name == "" {
		return "/apis/extensions/v1beta1/namespaces/" + namespace + "/replicasets"
	}
	return "/apis/extensions/v1beta1/namespaces/" + namespace + "/replicasets/" + name
}

// GetReplicaSet fetches a single ReplicaSet
func (c *Client) GetReplicaSet(namespace, name string) (*k8s.ReplicaSet, error) {
	var out k8s.ReplicaSet
	_, err := c.do("GET", replicasetGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ReplicaSet")
	}
	return &out, nil
}

// CreateReplicaSet creates a new ReplicaSet. This will fail if it already exists.
func (c *Client) CreateReplicaSet(namespace string, item *k8s.ReplicaSet) (*k8s.ReplicaSet, error) {
	item.TypeMeta.Kind = "ReplicaSet"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ReplicaSet
	_, err := c.do("POST", replicasetGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ReplicaSet")
	}
	return &out, nil
}

// ListReplicaSets lists all ReplicaSets in a namespace
func (c *Client) ListReplicaSets(namespace string, opts *k8s.ListOptions) (*k8s.ReplicaSetList, error) {
	var out k8s.ReplicaSetList
	_, err := c.do("GET", replicasetGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ReplicaSets")
	}
	return &out, nil
}

// DeleteReplicaSet deletes a single ReplicaSet. It will error if the ReplicaSet does not exist.
func (c *Client) DeleteReplicaSet(namespace, name string) error {
	_, err := c.do("DELETE", replicasetGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete ReplicaSet")
}

// UpdateReplicaSet will update in place a single ReplicaSet. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateReplicaSet(namespace string, item *k8s.ReplicaSet) (*k8s.ReplicaSet, error) {
	item.TypeMeta.Kind = "ReplicaSet"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ReplicaSet
	_, err := c.do("PUT", replicasetGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update ReplicaSet")
	}
	return &out, nil
}
