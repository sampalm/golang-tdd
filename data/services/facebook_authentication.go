package services

import (
	"errors"
	"tdd/data/contracts/api"
	"tdd/data/contracts/repos"
	domainErrs "tdd/domain/errors"
	"tdd/domain/features/facebook"
)

type FacebookAuthenticatonService struct {
	loadFacebookUser          api.LoadFacebookUserApi
	loadUserAccountRepo       repos.LoadUserAccountRepository
	createFacebookAccountRepo repos.CreateFacebookAccountRepository
}

func NewFacebookAuthenticatonService(lf api.LoadFacebookUserApi, lu repos.LoadUserAccountRepository, fb repos.CreateFacebookAccountRepository) FacebookAuthenticatonService {
	return FacebookAuthenticatonService{
		loadFacebookUser:          lf,
		loadUserAccountRepo:       lu,
		createFacebookAccountRepo: fb,
	}
}

func (fs FacebookAuthenticatonService) Perform(params facebook.Params) (facebook.Result, error) {
	res, err := fs.loadFacebookUser.Load(api.Params(params))
	if err != nil {
		return facebook.Result{}, domainErrs.AuthenticationError{Err: err}
	}

	_, err = fs.loadUserAccountRepo.Load(repos.LoadUserAccountParams{Email: res.User.Email})
	if err != nil {
		if errors.Is(err, domainErrs.ErrNoData) {
			userFb := repos.CreateFacebookAccountParams{
				Email:      res.User.Email,
				Name:       res.User.Name,
				FacebookID: res.User.FacebookID,
			}
			fs.createFacebookAccountRepo.CreateFromFacebook(userFb)
		}
	}
	return facebook.Result{}, domainErrs.AuthenticationError{}
}
