package http_test

import (
	"testing"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
)

func TestConfigMapList(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		list, err := c.ListConfigMaps(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
	})
}

func TestConfigMapCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		in := &client.ConfigMap{
			ObjectMeta: client.ObjectMeta{
				Name: "test-config",
			},
		}
		out, err := c.CreateConfigMap(n.Name, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		list, err := c.ListConfigMaps(n.Name, nil)
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Items) > 0, "should not be empty")

		out, err = c.GetConfigMap(n.Name, in.Name)
		assert.Nil(t, err)
		assert.NotNil(t, out)

		in.Data = map[string][]byte{
			"foo": []byte("value"),
		}

		out, err = c.UpdateConfigMap(n.Name, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)
		assert.True(t, len(out.Data) > 0, "should not be empty")

		err = c.DeleteConfigMap(n.Name, in.Name)
		assert.Nil(t, err)
	})
}
