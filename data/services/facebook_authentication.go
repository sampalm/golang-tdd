package services

import (
	"tdd/data/contracts/api"
	"tdd/domain/errors"
	"tdd/domain/features/facebook"
)

type FacebookAuthenticatonService struct {
	loadFacebookUser api.LoadFacebookUserApi
}

func NewFacebookAuthenticatonService(lf api.LoadFacebookUserApi) FacebookAuthenticatonService {
	return FacebookAuthenticatonService{
		loadFacebookUser: lf,
	}
}

func (fs FacebookAuthenticatonService) Perform(params facebook.Params) facebook.Result {
	fs.loadFacebookUser.LoadUser(api.Params(params))
	return facebook.Result{
		Err: &errors.AuthenticationError{},
	}
}
