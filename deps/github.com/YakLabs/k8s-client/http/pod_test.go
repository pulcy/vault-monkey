package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestPodList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListPods(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

func TestPodCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		/*
			in := &client.Pod{
				ObjectMeta: client.ObjectMeta{
					Name: "test-pod",
				},
			}
			out, err := c.CreatePod(n.Name, in)
			assert.Nil(t, err)
			assert.NotNil(t, out)

			list, err := c.ListPods(n.Name, nil)
			assert.Nil(t, err)
			assert.NotNil(t, list)
			assert.True(t, len(list.Items) > 0, "should not be empty")

			out, err = c.GetPod(n.Name, in.Name)
			assert.Nil(t, err)
			assert.NotNil(t, out)

			err = c.DeletePod(n.Name, in.Name)
			assert.Nil(t, err)
		*/
	})
}
