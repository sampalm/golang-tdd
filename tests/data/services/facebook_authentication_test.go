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

type MockCreateFacebookAccountProxy interface {
	// api interface
	repos.CreateFacebookAccountRepository

	// mock methods
	MockProxy
}

// SutTypes provides all the interfaces required to test facebook authentication feature.
type SutTypes struct {
	service                  services.FacebookAuthenticatonService
	loadFacebookUserApiSpy   MockApiProxy
	loadUserAccountRepoSpy   MockRepositoryProxy
	createFacebookAccountSpy MockCreateFacebookAccountProxy
}

func makeSut() SutTypes {
	loadFacebookUserApiSpy := new(LoadFacebookUserApiSpy)
	loadUserAccountRepoSpy := new(LoadUserAccountRepositorySpy)
	createFacebookAccountSpy := new(CreateFacebookAccountSpy)
	return SutTypes{
		service:                  services.NewFacebookAuthenticatonService(loadFacebookUserApiSpy, loadUserAccountRepoSpy, createFacebookAccountSpy),
		loadFacebookUserApiSpy:   loadFacebookUserApiSpy,
		loadUserAccountRepoSpy:   loadUserAccountRepoSpy,
		createFacebookAccountSpy: createFacebookAccountSpy,
	}
}

// LoadFacebookUserApiSpy is a mock implementation to test api.LoadFacebookUserApi
type LoadFacebookUserApiSpy struct {
	mock LoadFacebookUserApiSpyMock
}

func (lfs *LoadFacebookUserApiSpy) Load(params api.Params) (api.Result, error) {
	lfs.mock.token = params.Token
	lfs.mock.callsCount.Add(1)

	return lfs.mock.resultApi, lfs.mock.errApi
}

// LoadFacebookUserApiSpyMock allows to inject and retrive mocked data from LoadFacebookUserApiSpy implementation.
type LoadFacebookUserApiSpyMock struct {
	token      string
	resultApi  api.Result
	resultRepo repos.LoadUserAccountResult
	resultFb   repos.CreateFacebookAccountResult
	errApi     error
	errRepo    error
	errFb      error
	callsCount atomic.Int64
}

func (lfs *LoadFacebookUserApiSpy) Set(params LoadFacebookUserApiSpyMock) {
	lfs.mock.token = params.token
	lfs.mock.resultApi = params.resultApi
	lfs.mock.errApi = params.errApi
}
func (lfs *LoadFacebookUserApiSpy) Unset() {
	lfs.mock = LoadFacebookUserApiSpyMock{}
}
func (lfs *LoadFacebookUserApiSpy) Retrive() LoadFacebookUserApiSpyMock {
	return lfs.mock
}

func (lfs *LoadFacebookUserApiSpy) HaveBeenCalledTimes(n int) error {
	if n != int(lfs.mock.callsCount.Load()) {
		return fmt.Errorf("LoadFacebookUserApiSpy: expected to be called %d - has been called %d", n, int(lfs.mock.callsCount.Load()))
	}
	return nil
}

type LoadUserAccountRepositorySpy struct {
	mock LoadFacebookUserApiSpyMock
}

func (lfs *LoadUserAccountRepositorySpy) Load(params repos.LoadUserAccountParams) (repos.LoadUserAccountResult, error) {
	lfs.mock.callsCount.Add(1)
	return lfs.mock.resultRepo, lfs.mock.errRepo
}

func (lfs *LoadUserAccountRepositorySpy) Set(params LoadFacebookUserApiSpyMock) {
	lfs.mock.resultRepo = params.resultRepo
	lfs.mock.errRepo = params.errRepo
}
func (lfs *LoadUserAccountRepositorySpy) Unset() {
	lfs.mock = LoadFacebookUserApiSpyMock{}
}
func (lfs *LoadUserAccountRepositorySpy) Retrive() LoadFacebookUserApiSpyMock {
	return lfs.mock
}

func (lfs *LoadUserAccountRepositorySpy) HaveBeenCalledTimes(n int) error {
	if n != int(lfs.mock.callsCount.Load()) {
		return fmt.Errorf("LoadUserAccountRepositorySpy: expected to be called %d - has been called %d", n, int(lfs.mock.callsCount.Load()))
	}
	return nil
}

type CreateFacebookAccountSpy struct {
	mock LoadFacebookUserApiSpyMock
}

func (lfs *CreateFacebookAccountSpy) CreateFromFacebook(params repos.CreateFacebookAccountParams) error {
	lfs.mock.callsCount.Add(1)
	return lfs.mock.errFb
}

func (lfs *CreateFacebookAccountSpy) Set(params LoadFacebookUserApiSpyMock) {
	lfs.mock.resultFb = params.resultFb
	lfs.mock.errFb = params.errFb
}
func (lfs *CreateFacebookAccountSpy) Unset() {
	lfs.mock = LoadFacebookUserApiSpyMock{}
}
func (lfs *CreateFacebookAccountSpy) Retrive() LoadFacebookUserApiSpyMock {
	return lfs.mock
}

func (lfs *CreateFacebookAccountSpy) HaveBeenCalledTimes(n int) error {
	if n != int(lfs.mock.callsCount.Load()) {
		return fmt.Errorf("CreateFacebookAccountSpy: expected to be called %d - has been called %d", n, int(lfs.mock.callsCount.Load()))
	}
	return nil
}

//
// FACEBOOK AUTHENTICATION SERVICE
// UNIT TESTS
//

func TestFacebookAuthenticatonService(t *testing.T) {
	const token = "any_token"
	var defaultResultApi = api.Result{User: &api.FacebookUser{
		Name:       "any_fb_name",
		Email:      "any_fb_email",
		FacebookID: "any_fb_id",
	}}

	t.Run("should call LoadFacebookUserApi with correct params", func(t *testing.T) {
		assert.New(t)
		sut := makeSut()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{resultApi: defaultResultApi})

		sut.service.Perform(facebook.Params{Token: token})
		assert.Equal(t, token, sut.loadFacebookUserApiSpy.Retrive().token)

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))
	})

	t.Run("should return AuthenticationError when LoadFacebookUserApi returns no data", func(t *testing.T) {
		assert.New(t)
		sut := makeSut()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{errApi: domainErrs.ErrNotAuthenticated})

		_, err := sut.service.Perform(facebook.Params{Token: token})
		assert.ErrorIs(t, domainErrs.AuthenticationError{Err: domainErrs.ErrNotAuthenticated}, err)

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))

	})

	t.Run("should call LoadUserByEmailRep when LoadFacebookUserApi returns data", func(t *testing.T) {
		assert.New(t)
		sut := makeSut()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{resultApi: defaultResultApi, errApi: nil})

		_, err := sut.service.Perform(facebook.Params{Token: token})
		assert.ErrorIs(t, err, domainErrs.AuthenticationError{})

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))
		assert.NoError(t, sut.loadUserAccountRepoSpy.HaveBeenCalledTimes(1))

	})

	t.Run("should call CreateUserByEmailRep when LoadFacebookUserApi returns no data", func(t *testing.T) {
		assert.New(t)
		sut := makeSut()
		sut.loadFacebookUserApiSpy.Set(LoadFacebookUserApiSpyMock{resultApi: defaultResultApi, errApi: nil})
		sut.loadUserAccountRepoSpy.Set(LoadFacebookUserApiSpyMock{errRepo: domainErrs.ErrNoData})

		_, err := sut.service.Perform(facebook.Params{Token: token})
		assert.ErrorIs(t, domainErrs.AuthenticationError{}, err)

		assert.NoError(t, sut.loadFacebookUserApiSpy.HaveBeenCalledTimes(1))
		assert.NoError(t, sut.loadUserAccountRepoSpy.HaveBeenCalledTimes(1))
		assert.NoError(t, sut.createFacebookAccountSpy.HaveBeenCalledTimes(1))

	})

}
