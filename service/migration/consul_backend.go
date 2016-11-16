package migration

import (
	"net/url"
	"strings"

	"github.com/hashicorp/consul/api"
)

type consulBackend struct {
	path   string
	client *api.Client
	kv     *api.KV
}

func NewConsulBackend(address string) (Backend, error) {
	url, err := url.Parse(address)
	if err != nil {
		return nil, maskAny(err)
	}
	path := url.Path

	// Ensure path is suffixed but not prefixed
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	consulConf := api.DefaultConfig()
	consulConf.Address = url.Host

	client, err := api.NewClient(consulConf)
	if err != nil {
		return nil, maskAny(err)
	}

	return &consulBackend{
		path:   path,
		client: client,
		kv:     client.KV(),
	}, nil
}

func (b *consulBackend) Get(key string) ([]byte, error) {
	pair, _, err := b.kv.Get(b.path+key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, nil
	}
	return pair.Value, nil
}

func (b *consulBackend) Set(key string, value []byte) error {
	pair := &api.KVPair{
		Key:   b.path + key,
		Value: value,
	}

	_, err := b.kv.Put(pair, nil)
	return err
}

func (b *consulBackend) List(key string) ([]string, error) {
	scan := b.path + key

	// The TrimPrefix call below will not work correctly if we have "//" at the
	// end. This can happen in cases where you are e.g. listing the root of a
	// prefix in a logical backend via "/" instead of ""
	if strings.HasSuffix(scan, "//") {
		scan = scan[:len(scan)-1]
	}

	out, _, err := b.kv.Keys(scan, "/", nil)
	for idx, val := range out {
		out[idx] = strings.TrimPrefix(val, scan)
	}

	return out, err
}
