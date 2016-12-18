package client

import "github.com/pkg/errors"

// IsNotFoundError can be used to check if the error was a not found error.
func IsNotFoundError(err error) bool {
	s, ok := errors.Cause(err).(*Status)
	return ok && s.Code == 404
}
