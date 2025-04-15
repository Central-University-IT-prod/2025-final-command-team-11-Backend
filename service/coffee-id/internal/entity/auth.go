package entity

type Tokens struct {
	Access  string
	Refresh string
}

type AccessData struct {
	Token string
	Roles []string
}
