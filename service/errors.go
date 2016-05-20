// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/juju/errgo"
)

var (
	InvalidArgumentError = errgo.New("invalid argument")
	VaultError           = errgo.New("vault error")
	SecretNotFoundError  = errgo.New("secret not found")
	maskAny              = errgo.MaskFunc(errgo.Any)
)

func IsVault(err error) bool {
	return errgo.Cause(err) == VaultError
}

func IsSecretNotFound(err error) bool {
	return errgo.Cause(err) == SecretNotFoundError
}

type AggregateError struct {
	errors []error
}

func collectErrorsFromChannel(errors chan error) error {
	ae := &AggregateError{}
	for err := range errors {
		if err != nil {
			ae.errors = append(ae.errors, err)
		}
	}
	switch len(ae.errors) {
	case 0:
		return nil // no error
	case 1:
		return ae.errors[0]
	default:
		return ae
	}
}

func (ae *AggregateError) Error() string {
	l := []string{}
	for _, err := range ae.errors {
		l = append(l, err.Error())
	}
	return strings.Join(l, ", ")
}

func Describe(err error) string {
	if urlErr, ok := err.(*url.Error); ok {
		return fmt.Sprintf("Op=%s, URL=%s, Error=%s (%#v)", urlErr.Op, urlErr.URL, urlErr.Err.Error(), err)
	}
	return err.Error()
}
