package api

type LoadFacebookUserApi interface {
	Load(Params) (Result, error)
}

type Params struct {
	Token string
}

type Result struct {
	User FacebookUser
}

type FacebookUser struct {
	FacebookID string
	Email      string
	Name       string
}
