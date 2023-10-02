package services

import (
	"errors"
	"tdd/data/contracts/api"
	"tdd/data/contracts/repos"
	domainErrs "tdd/domain/errors"
	"tdd/domain/features/facebook"
)

type userAccountRepo struct {
	repos.LoadUserAccountRepository
	repos.CreateFacebookAccountRepository
}

type FacebookAuthenticatonService struct {
	loadFacebookUser api.LoadFacebookUserApi
	userAccountRepo  userAccountRepo
}

func NewFacebookAuthenticatonService(lf api.LoadFacebookUserApi, lu repos.LoadUserAccountRepository, fb repos.CreateFacebookAccountRepository) FacebookAuthenticatonService {
	return FacebookAuthenticatonService{
		loadFacebookUser: lf,
		userAccountRepo:  userAccountRepo{lu, fb},
	}
}

func (fs FacebookAuthenticatonService) Perform(params facebook.Params) (facebook.Result, error) {
	res, err := fs.loadFacebookUser.Load(api.Params(params))
	if err != nil {
		return facebook.Result{}, domainErrs.AuthenticationError{Err: err}
	}

	_, err = fs.userAccountRepo.Load(repos.LoadUserAccountParams{Email: res.User.Email})
	if err != nil {
		if errors.Is(err, domainErrs.ErrNoData) {
			userFb := repos.CreateFacebookAccountParams{
				Email:      res.User.Email,
				Name:       res.User.Name,
				FacebookID: res.User.FacebookID,
			}
			fs.userAccountRepo.CreateFromFacebook(userFb)
		}
	}
	return facebook.Result{}, domainErrs.AuthenticationError{}
}
