package services_test

import (
	"sync/atomic"
	"tdd/data/contracts/api"
	"tdd/data/services"
	domainErrs "tdd/domain/errors"
	"tdd/domain/features/facebook"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LoadFacebookUserApiSpy struct {
	token      string
	callsCount atomic.Int64
	result     api.Result
}

func (lfs *LoadFacebookUserApiSpy) LoadUser(params api.Params) api.Result {
	lfs.token = params.Token
	lfs.callsCount.Add(1)
	return lfs.result
}

func TestFacebookAuthenticatonService(t *testing.T) {

	t.Run("should call LoadFacebookUserApi with correct params", func(t *testing.T) {
		var loadFacebookUserApi = &LoadFacebookUserApiSpy{}
		var sut = services.NewFacebookAuthenticatonService(loadFacebookUserApi)

		sut.Perform(facebook.Params{Token: "any_token"})
		assert.Equal(t, "any_token", loadFacebookUserApi.token)
		assert.Equal(t, int64(1), loadFacebookUserApi.callsCount.Load())
	})

	t.Run("should return AuthenticationError when LoadFacebookUserApi returns nil", func(t *testing.T) {
		var loadFacebookUserApi = &LoadFacebookUserApiSpy{}
		loadFacebookUserApi.result = api.Result{User: nil}
		var sut = services.NewFacebookAuthenticatonService(loadFacebookUserApi)

		result := sut.Perform(facebook.Params{Token: "any_token"})
		assert.ErrorIs(t, domainErrs.AuthenticationError{}, *result.Err)
	})

}
