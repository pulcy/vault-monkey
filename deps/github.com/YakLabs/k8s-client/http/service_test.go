package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestServiceList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListServices(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

/*
func TestServiceCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {

		in := &client.Service{
			ObjectMeta: client.ObjectMeta{
				Name: "test-service",
			},
		}
		out, err := c.CreateService(n.Name, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		list, err := c.ListServices(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Items) > 0, "should not be empty")

		out, err = c.GetService(n.Name, in.Name)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		err = c.DeleteService(n.Name, in.Name)
		assert.Nil(t, err)
	})
}
*/
