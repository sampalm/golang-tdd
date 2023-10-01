package repos

type LoadUserAccountRepository interface {
	Load(params LoadUserAccountParams)
}
type LoadUserAccountParams struct {
	Email string
}
