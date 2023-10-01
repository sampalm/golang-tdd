package services

import (
	"tdd/data/contracts/api"
	"tdd/data/contracts/repos"
	"tdd/domain/errors"
	"tdd/domain/features/facebook"
)

type FacebookAuthenticatonService struct {
	loadFacebookUser    api.LoadFacebookUserApi
	loadUserAccountRepo repos.LoadUserAccountRepository
}

func NewFacebookAuthenticatonService(lf api.LoadFacebookUserApi, lu repos.LoadUserAccountRepository) FacebookAuthenticatonService {
	return FacebookAuthenticatonService{
		loadFacebookUser:    lf,
		loadUserAccountRepo: lu,
	}
}

func (fs FacebookAuthenticatonService) Perform(params facebook.Params) (facebook.Result, error) {
	res, err := fs.loadFacebookUser.Load(api.Params(params))
	if err != nil {
		return facebook.Result{}, errors.AuthenticationError{Err: err}
	}

	fs.loadUserAccountRepo.Load(repos.LoadUserAccountParams{Email: res.User.Email})
	return facebook.Result{}, errors.AuthenticationError{}
}
