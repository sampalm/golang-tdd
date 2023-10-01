package services_test

import (
	"tdd/domain/features/facebook"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FacebookAuthenticatonService struct {
	loadFacebookUser LoadFacebookUserApi
}

func NewFacebookAuthenticatonService(lf LoadFacebookUserApi) FacebookAuthenticatonService {
	return FacebookAuthenticatonService{
		loadFacebookUser: lf,
	}
}

func (fs FacebookAuthenticatonService) Perform(params facebook.Params) facebook.Result {
	fs.loadFacebookUser.LoadUser(Params(params))
	return facebook.Result{}
}

type LoadFacebookUserApi interface {
	LoadUser(Params)
}

type Params struct {
	Token string
}

type LoadFacebookUserApiSpy struct {
	token string
}

func (lfs *LoadFacebookUserApiSpy) LoadUser(params Params) {
	lfs.token = params.Token
}

func TestFacebookAuthenticatonService(t *testing.T) {
	var loadFacebookUserApi = &LoadFacebookUserApiSpy{}
	var sut = NewFacebookAuthenticatonService(loadFacebookUserApi)

	sut.Perform(facebook.Params{Token: "any_token"})
	assert.Equal(t, "any_token", loadFacebookUserApi.token)
}
