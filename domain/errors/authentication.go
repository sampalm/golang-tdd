package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNoData           = errors.New("no data")
	ErrInvalidParams    = errors.New("invalid params")
	ErrNotAuthenticated = errors.New("not authenticated")
)

type AuthenticationError struct {
	Err error
}

func (r AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %v", r.Err)
}
