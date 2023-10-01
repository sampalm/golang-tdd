package facebook

import (
	"tdd/domain/errors"
	"tdd/domain/models"
)

type FacebookAuthenticaton interface {
	Perform(params Params) Result
}

// Params is implemented on top of Command Pattern
type Params struct {
	Token string
}

// Result is implemented on top of Command Pattern
type Result struct {
	AccessToken models.AccessToken
	Err         *errors.AuthenticationError
}
