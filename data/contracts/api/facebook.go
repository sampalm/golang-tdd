package api

type LoadFacebookUserApi interface {
	LoadUser(Params) Result
}

type Params struct {
	Token string
}

type Result struct {
	User *string
}
