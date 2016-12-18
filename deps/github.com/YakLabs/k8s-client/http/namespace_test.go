package http_test

import (
	"testing"
	"time"

	"github.com/YakLabs/k8s-client"
	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//helper for using a test namespace
func withTestNamespace(t *testing.T, f func(*testing.T, *http.Client, *client.Namespace)) {
	c := testClient(t)

	for {
		// deletion is not immediate, so wait until its gone
		_, err := c.GetNamespace("test123")
		// should make sure this is a not found error, but it'll fail later
		// if its something else
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	ns := client.Namespace{
		ObjectMeta: client.ObjectMeta{
			Name: "test123",
		},
	}
	out, err := c.CreateNamespace(&ns)
	require.Nil(t, err)
	require.NotNil(t, out)

	f(t, c, out)
	err = c.DeleteNamespace("test123")
	require.Nil(t, err)
}

func TestNamespaceGet(t *testing.T) {
	c := testClient(t)
	n, err := c.GetNamespace("default")
	assert.Nil(t, err)
	assert.NotNil(t, n)
	assert.Equal(t, "default", n.Name, "name should be equal")
}

func TestNamespaceList(t *testing.T) {
	c := testClient(t)
	list, err := c.ListNamespaces(nil)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.True(t, len(list.Items) > 0, "list should not be empty")
}

func TestNamespaceCreate(t *testing.T) {
	withTestNamespace(t, func(t *testing.T, c *http.Client, n *client.Namespace) {
		n.ObjectMeta.Labels = map[string]string{
			"foo": "bar",
		}
		out, err := c.UpdateNamespace(n)
		assert.Nil(t, err)
		assert.NotNil(t, out)
		assert.NotNil(t, out.ObjectMeta.Labels)
		assert.Equal(t, "bar", out.ObjectMeta.Labels["foo"], "labels should be equal")
	})
}
