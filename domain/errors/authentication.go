package errors

import "fmt"

type AuthenticationError struct {
	Err error
}

func (r AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %v", r.Err)
}
