package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestServiceAccountList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListServiceAccounts(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

func TestServiceAccountCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		in := &client.ServiceAccount{
			ObjectMeta: client.ObjectMeta{
				Name: "test-service-account",
			},
		}
		out, err := c.CreateServiceAccount(n.Name, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		list, err := c.ListServiceAccounts(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Items) > 0, "should not be empty")

		out, err = c.GetServiceAccount(n.Name, in.Name)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		err = c.DeleteServiceAccount(n.Name, in.Name)
		assert.Nil(t, err)
	})
}
