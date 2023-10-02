package repos

type LoadUserAccountRepository interface {
	Load(params LoadUserAccountParams) (LoadUserAccountResult, error)
}
type LoadUserAccountParams struct {
	Email string
}
type LoadUserAccountResult struct {
	User *UserAccount
}

type UserAccount struct {
}

type CreateFacebookAccountRepository interface {
	CreateFromFacebook(params CreateFacebookAccountParams) error
}
type CreateFacebookAccountParams struct {
	Email      string
	Name       string
	FacebookID string
}
type CreateFacebookAccountResult struct {
	User *UserAccount
}
