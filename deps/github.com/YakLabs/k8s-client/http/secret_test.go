package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestSecretList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListSecrets(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

func TestSecretCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		in := client.NewSecret(n.Name, "test-secret")

		out, err := c.CreateSecret(n.Name, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		list, err := c.ListSecrets(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Items) > 0, "should not be empty")

		out, err = c.GetSecret(n.Name, in.Name)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		err = c.DeleteSecret(n.Name, in.Name)
		assert.Nil(t, err)
	})
}
