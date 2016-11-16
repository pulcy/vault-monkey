package migration

import (
	"context"
	"encoding/base64"
	"net/url"
	"path"
	"strings"

	"github.com/coreos/etcd/client"
)

const (
	EtcdNodeFilePrefix = "."
)

type etcdBackend struct {
	path string
	kAPI client.KeysAPI
}

func NewEtcdBackend(address string) (Backend, error) {
	url, err := url.Parse(address)
	if err != nil {
		return nil, maskAny(err)
	}
	path := url.Path

	// Ensure path is prefixed.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	url.Path = ""
	endpoint := url.String()

	c, err := client.New(client.Config{
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, err
	}
	kAPI := client.NewKeysAPI(c)

	return &etcdBackend{
		path: path,
		kAPI: kAPI,
	}, nil
}

func (b *etcdBackend) Get(key string) ([]byte, error) {
	getOpts := &client.GetOptions{
		Recursive: false,
		Sort:      false,
	}
	response, err := b.kAPI.Get(context.Background(), b.nodePath(key), getOpts)
	if err != nil {
		if errorIsMissingKey(err) {
			return nil, nil
		}
		return nil, err
	}

	// Decode the stored value from base-64.
	value, err := base64.StdEncoding.DecodeString(response.Node.Value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (b *etcdBackend) Set(key string, value []byte) error {
	encValue := base64.StdEncoding.EncodeToString(value)

	_, err := b.kAPI.Set(context.Background(), b.nodePath(key), encValue, nil)
	return err

}

func (b *etcdBackend) List(key string) ([]string, error) {
	// Set a directory path from the given prefix.
	path := b.nodePathDir(key)

	// Get the directory, non-recursively, from etcd. If the directory is
	// missing, we just return an empty list of contents.
	getOpts := &client.GetOptions{
		Recursive: false,
		Sort:      true,
	}
	response, err := b.kAPI.Get(context.Background(), path, getOpts)
	if err != nil {
		if errorIsMissingKey(err) {
			return []string{}, nil
		}
		return nil, err
	}

	out := make([]string, len(response.Node.Nodes))
	for i, node := range response.Node.Nodes {

		// etcd keys include the full path, so let's trim the prefix directory
		// path.
		name := strings.TrimPrefix(node.Key, path)

		// Check if this node is itself a directory. If it is, add a trailing
		// slash; if it isn't remove the node file prefix.
		if node.Dir {
			out[i] = name + "/"
		} else {
			out[i] = name[1:]
		}
	}
	return out, nil
}

// nodePath returns an etcd filepath based on the given key.
func (b *etcdBackend) nodePath(key string) string {
	return path.Join(b.path, path.Dir(key), EtcdNodeFilePrefix+path.Base(key))
}

// nodePathDir returns an etcd directory path based on the given key.
func (b *etcdBackend) nodePathDir(key string) string {
	return path.Join(b.path, key) + "/"
}

// errorIsMissingKey returns true if the given error is an etcd error with an
// error code corresponding to a missing key.
func errorIsMissingKey(err error) bool {
	etcdErr, ok := err.(client.Error)
	return ok && etcdErr.Code == client.ErrorCodeKeyNotFound
}
