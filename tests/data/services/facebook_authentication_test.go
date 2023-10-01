package services_test

import (
	"fmt"
	"sync/atomic"
	"tdd/data/contracts/api"
	"tdd/data/contracts/repos"
	"tdd/data/services"
	domainErrs "tdd/domain/errors"
	"tdd/domain/features/facebook"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockProxy provides all contracts required to decouple all dependencies from loadFacebookUserApiSpy
type MockProxy interface {
	// mock methods
	Set(params LoadFacebookUserApiSpyMock)
	Retrive() LoadFacebookUserApiSpyMock
	Unset()
	HaveBeenCalledTimes(int) error
}
type MockApiProxy interface {
	// api interface
	api.LoadFacebookUserApi

	// mock methods
	MockProxy
}

type MockRepositoryProxy interface {
	// api interface
	repos.LoadUserAccountRepository

	// mock methods
	MockProxy
}

// SutTypes provides all the interfaces required to test facebook authentication feature.
type SutTypes struct {
	service                services.FacebookAuthenticatonService
	loadFacebookUserApiSpy MockApiProxy
	loadUserAccountRepoSpy MockRepositoryProxy
}

func makeSut() SutTypes {
	loadFacebookUserApiSpy := new(LoadFacebookUserApiSpy)
	loadUserAccountRepoSpy := new(LoadUserAccountRepositorySpy)
	return SutTypes{
		service:                services.NewFacebookAuthenticatonService(loadFacebookUserApiSpy, loadUserAccountRepoSpy),
		loadFacebookUserApiSpy: loadFacebookUserApiSpy,
		loadUserAccountRepoSpy: loadUserAccountRepoSpy,
	}
}

// LoadFacebookUserApiSpy is a mock implementation to test api.LoadFacebookUserApi
type LoadFacebookUserApiSpy struct {
	mock LoadFacebookUserApiSpyMock
}

func (lfs *LoadFacebookUserApiSpy) Load(params api.Params) (api.Result, error) {
	lfs.mock.token = params.Token
	lfs.mock.callsCount.Add(1)
	return lfs.mock.result, nil
}

// LoadFacebookUserApiSpyMock allows to inject and retrive mocked data from LoadFacebookUserApiSpy implementation.
type LoadFacebookUserApiSpyMock struct {
	token      string
	result     api.Result
	callsCount atomic.Int64
}

func (lfs *LoadFacebookUserApiSpy) Set(params LoadFacebookUserApiSpyMock) {
	lfs.mock.token = params.token
}
func (lfs *LoadFacebookUserApiSpy) Unset() {
	lfs.mock = LoadFacebookUserApiSpyMock{}
}
func (lfs *LoadFacebookUserApiSpy) Retrive() LoadFacebookUserApiSpyMock {
	return lfs.mock
}

func (lfs *LoadFacebookUserApiSpy) HaveBeenCalledTimes(n int) error {
	if n != int(lfs.mock.callsCount.Load()) {
		return fmt.Errorf("expected to be called %d - has been called %d", n, int(lfs.mock.callsCount.Load()))
	}
	return nil
}

type LoadUserAccountRepositorySpy struct {
	mock LoadFacebookUserApiSpyMock
}

func (lfs *LoadUserAccountRepositorySpy) Load(params repos.LoadUserAccountParams) {
	lfs.mock.callsCount.Add(1)
}

func (lfs *LoadUserAccountRepositorySpy) Set(params LoadFacebookUserApiSpyMock) {
	lfs.mock.token = params.token
}
func (lfs *LoadUserAccountRepositorySpy) Unset() {
	lfs.mock = LoadFacebookUserApiSpyMock{}
}
func (lfs *LoadUserAccountRepositorySpy) Retrive() LoadFacebookUserApiSpyMock {
	return lfs.mock
}

func (lfs *LoadUserAccountRepositorySpy) HaveBeenCalledTimes(n int) error {
	if n != int(lfs.mock.callsCount.Load()) {
		return fmt.Errorf("expected to be called %d - has been called %d", n, int(lfs.mock.callsCount.Load()))
	}
	return nil
}

//
// FACEBOOK AUTHENTICATION SERVICE
// UNIT TESTS
//

func TestFacebookAuthenticatonService(t *testing.T) {
	const token = "any_token"

	t.Run("should call LoadFacebookUserApi with correct params", func(t *testing.T) {
		assert.New(t)
		sut := makeSut()
		defer sut.loadFacebookUserApiSpy.Unset()

		sut.service.Perform(facebook.Params{Token: token})
		assert.Equal(t, token, sut.loadFacebookUserApiSpy.Retrive().token)

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))
	})

	t.Run("should return AuthenticationError when LoadFacebookUserApi returns nil", func(t *testing.T) {
		assert.New(t)
		expectResult := api.Result{}
		sut := makeSut()
		defer sut.loadFacebookUserApiSpy.Unset()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{result: expectResult})

		_, err := sut.service.Perform(facebook.Params{Token: token})
		assert.ErrorIs(t, domainErrs.AuthenticationError{}, err)

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))

	})

	t.Run("should call LoadUserByEmailRep when LoadFacebookUserApi returns data", func(t *testing.T) {
		assert.New(t)
		expectResult := api.Result{User: api.FacebookUser{
			Name:       "any_fb_name",
			Email:      "any_fb_email",
			FacebookID: "any_fb_id",
		}}
		sut := makeSut()
		defer sut.loadFacebookUserApiSpy.Unset()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{result: expectResult})

		// assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))

		_, err := sut.service.Perform(facebook.Params{Token: token})
		assert.ErrorIs(t, domainErrs.AuthenticationError{}, err)

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))
		assert.NoError(t, sut.loadUserAccountRepoSpy.HaveBeenCalledTimes(1))

	})

}
