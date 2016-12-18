package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func horizontalpodautoscalerGeneratePath(namespace, name string) string {
	if name == "" {
		return "/apis/autoscaling/v1/namespaces/" + namespace + "/horizontalpodautoscalers"
	}
	return "/apis/autoscaling/v1/namespaces/" + namespace + "/horizontalpodautoscalers/" + name
}

// GetHorizontalPodAutoscaler fetches a single HorizontalPodAutoscaler
func (c *Client) GetHorizontalPodAutoscaler(namespace, name string) (*k8s.HorizontalPodAutoscaler, error) {
	var out k8s.HorizontalPodAutoscaler
	_, err := c.do("GET", horizontalpodautoscalerGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get HorizontalPodAutoscaler")
	}
	return &out, nil
}

// CreateHorizontalPodAutoscaler creates a new HorizontalPodAutoscaler. This will fail if it already exists.
func (c *Client) CreateHorizontalPodAutoscaler(namespace string, item *k8s.HorizontalPodAutoscaler) (*k8s.HorizontalPodAutoscaler, error) {
	item.TypeMeta.Kind = "HorizontalPodAutoscaler"
	item.TypeMeta.APIVersion = "autoscaling/v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.HorizontalPodAutoscaler
	_, err := c.do("POST", horizontalpodautoscalerGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HorizontalPodAutoscaler")
	}
	return &out, nil
}

// ListHorizontalPodAutoscalers lists all HorizontalPodAutoscalers in a namespace
func (c *Client) ListHorizontalPodAutoscalers(namespace string, opts *k8s.ListOptions) (*k8s.HorizontalPodAutoscalerList, error) {
	var out k8s.HorizontalPodAutoscalerList
	_, err := c.do("GET", horizontalpodautoscalerGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list HorizontalPodAutoscalers")
	}
	return &out, nil
}

// DeleteHorizontalPodAutoscaler deletes a single HorizontalPodAutoscaler. It will error if the HorizontalPodAutoscaler does not exist.
func (c *Client) DeleteHorizontalPodAutoscaler(namespace, name string) error {
	_, err := c.do("DELETE", horizontalpodautoscalerGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete HorizontalPodAutoscaler")
}

// UpdateHorizontalPodAutoscaler will update in place a single HorizontalPodAutoscaler. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateHorizontalPodAutoscaler(namespace string, item *k8s.HorizontalPodAutoscaler) (*k8s.HorizontalPodAutoscaler, error) {
	item.TypeMeta.Kind = "HorizontalPodAutoscaler"
	item.TypeMeta.APIVersion = "autoscaling/v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.HorizontalPodAutoscaler
	_, err := c.do("PUT", horizontalpodautoscalerGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update HorizontalPodAutoscaler")
	}
	return &out, nil
}
