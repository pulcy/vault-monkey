package http_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeList(t *testing.T) {
	c := testClient(t)
	list, err := c.ListNodes(nil)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.True(t, len(list.Items) > 0, "list should not be empty")
}

func TestNodeGet(t *testing.T) {
	c := testClient(t)
	list, err := c.ListNodes(nil)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.True(t, len(list.Items) > 0, "list should not be empty")

	out, err := c.GetNode(list.Items[0].Name)
	assert.Nil(t, err)
	assert.NotNil(t, out)
}
