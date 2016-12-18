package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestHorizontalPodAutoscalerList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListHorizontalPodAutoscalers(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

func TestHorizontalPodAutoscalerCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		in := client.NewHorizontalPodAutoscaler(n.Name, "test-hpa")
		in.Spec.MaxReplicas = 3
		in.Spec.ScaleTargetRef = client.CrossVersionObjectReference{
			Kind:       "ReplicationController",
			Name:       "test",
			APIVersion: "v1",
		}

		out, err := c.CreateHorizontalPodAutoscaler(n.Name, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		list, err := c.ListHorizontalPodAutoscalers(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Items) > 0, "should not be empty")

		out, err = c.GetHorizontalPodAutoscaler(n.Name, in.Name)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		err = c.DeleteHorizontalPodAutoscaler(n.Name, in.Name)
		assert.Nil(t, err)

	})
}
