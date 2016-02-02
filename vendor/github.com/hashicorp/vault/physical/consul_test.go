package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

func TestConsulBackend(t *testing.T) {
	addr := os.Getenv("CONSUL_ADDR")
	if addr == "" {
		t.SkipNow()
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	b, err := NewBackend("consul", map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "256",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}

func TestConsulHABackend(t *testing.T) {
	addr := os.Getenv("CONSUL_ADDR")
	if addr == "" {
		t.SkipNow()
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	b, err := NewBackend("consul", map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "-1",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ha, ok := b.(HABackend)
	if !ok {
		t.Fatalf("consul does not implement HABackend")
	}
	testHABackend(t, ha, ha)

	detect, ok := b.(AdvertiseDetect)
	if !ok {
		t.Fatalf("consul does not implement AdvertiseDetect")
	}
	host, err := detect.DetectHostAddr()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if host == "" {
		t.Fatalf("bad addr: %v", host)
	}
}
