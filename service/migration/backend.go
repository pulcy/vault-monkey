package migration

import "github.com/juju/errgo"

type Backend interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	List(key string) ([]string, error)
}

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)
