package service

import (
	"github.com/juju/errgo"
)

var (
	InvalidArgumentError = errgo.New("invalid argument")
	VaultError           = errgo.New("vault error")
	maskAny              = errgo.MaskFunc(errgo.Any)
)
