package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

// GetNode gets a single node.
func (c *Client) GetNode(name string) (*k8s.Node, error) {
	var out k8s.Node
	_, err := c.do("GET", "/api/v1/nodes/"+name, nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node")
	}
	return &out, nil
}

// CreateNode creates a single node. It will fail if it already exists.
func (c *Client) CreateNode(item *k8s.Node) (*k8s.Node, error) {
	item.TypeMeta.Kind = "Node"
	item.TypeMeta.APIVersion = "v1"

	var out k8s.Node
	_, err := c.do("POST", "/api/v1/nodes", item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create node")
	}
	return &out, nil
}

// ListNodes list all nodes, optionally filtering.
func (c *Client) ListNodes(opts *k8s.ListOptions) (*k8s.NodeList, error) {
	var out k8s.NodeList
	_, err := c.do("GET", "/api/v1/nodes?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list nodes")
	}
	return &out, nil
}

// DeleteNode removes a single node.
func (c *Client) DeleteNode(name string) error {
	_, err := c.do("DELETE", "/api/v1/nodes/"+name, nil, nil)
	return errors.Wrap(err, "failed to delete node")
}

// UpdateNode updates s sinle node.
func (c *Client) UpdateNode(item *k8s.Node) (*k8s.Node, error) {
	item.TypeMeta.Kind = "Node"
	item.TypeMeta.APIVersion = "v1"

	var out k8s.Node
	_, err := c.do("PUT", "/api/v1/nodes/"+item.Name, item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update node")
	}
	return &out, nil
}
