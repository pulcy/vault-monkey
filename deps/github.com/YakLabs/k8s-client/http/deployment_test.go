package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListDeployments(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

func TestDeploymentCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		/*
			in := &client.Deployment{
				ObjectMeta: client.ObjectMeta{
					Name: "test-deployment",
				},
			}
			out, err := c.CreateDeployment(n.Name, in)
			assert.Nil(t, err)
			assert.NotNil(t, out)

			list, err := c.ListDeployments(n.Name, nil)
			assert.Nil(t, err)
			assert.NotNil(t, list)
			assert.True(t, len(list.Items) > 0, "should not be empty")

			out, err = c.GetDeployment(n.Name, in.Name)
			assert.Nil(t, err)
			assert.NotNil(t, out)

			err = c.DeleteDeployment(n.Name, in.Name)
			assert.Nil(t, err)
		*/
	})
}
