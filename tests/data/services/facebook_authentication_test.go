package services_test

import (
	"sync/atomic"
	"tdd/data/contracts/api"
	"tdd/data/services"
	domainErrs "tdd/domain/errors"
	"tdd/domain/features/facebook"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProxy interface {
	// api interface
	api.LoadFacebookUserApi

	// mock methods
	Set(params LoadFacebookUserApiSpyMock)
	Retrive() LoadFacebookUserApiSpyMock

	// testify mock
	On(methodName string, arguments ...interface{}) *mock.Call
	AssertExpectations(mock.TestingT) bool
}

type SutTypes struct {
	service                services.FacebookAuthenticatonService
	loadFacebookUserApiSpy MockProxy
}

func makeSut() SutTypes {
	loadFacebookUserApiSpy := new(LoadFacebookUserApiSpy)
	return SutTypes{
		service:                services.NewFacebookAuthenticatonService(loadFacebookUserApiSpy),
		loadFacebookUserApiSpy: loadFacebookUserApiSpy,
	}
}

type LoadFacebookUserApiSpyMock struct {
	token      string
	result     api.Result
	callsCount atomic.Int64
}

type LoadFacebookUserApiSpy struct {
	mock LoadFacebookUserApiSpyMock
	mock.Mock
}

func (lfs *LoadFacebookUserApiSpy) Set(params LoadFacebookUserApiSpyMock) {
	lfs.mock.token = params.token
}
func (lfs *LoadFacebookUserApiSpy) Retrive() LoadFacebookUserApiSpyMock {
	return lfs.mock
}

func (lfs *LoadFacebookUserApiSpy) LoadUser(params api.Params) api.Result {
	lfs.mock.token = params.Token
	lfs.Called(params)
	lfs.mock.callsCount.Add(1)
	return lfs.mock.result
}

func TestFacebookAuthenticatonService(t *testing.T) {

	t.Run("should call LoadFacebookUserApi with correct params", func(t *testing.T) {
		assert.New(t)
		expectResult := api.Result{User: nil}
		sut := makeSut()

		sut.loadFacebookUserApiSpy.On("LoadUser", api.Params{Token: "any_token"}).Return(expectResult)

		sut.service.Perform(facebook.Params{Token: "any_token"})
		assert.Equal(t, "any_token", sut.loadFacebookUserApiSpy.Retrive().token)
		calls := sut.loadFacebookUserApiSpy.Retrive().callsCount
		assert.Equal(t, int64(1), calls.Load())

		sut.loadFacebookUserApiSpy.AssertExpectations(t)
	})

	t.Run("should return AuthenticationError when LoadFacebookUserApi returns nil", func(t *testing.T) {
		assert.New(t)
		expectResult := api.Result{User: nil}
		sut := makeSut()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{result: expectResult})

		sut.loadFacebookUserApiSpy.On("LoadUser", api.Params{Token: "any_token"}).Return(expectResult)

		result := sut.service.Perform(facebook.Params{Token: "any_token"})
		assert.ErrorIs(t, domainErrs.AuthenticationError{}, *result.Err)

		sut.loadFacebookUserApiSpy.AssertExpectations(t)

	})

}
