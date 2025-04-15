package entity

import "encoding/json"

type OAuthSession struct {
	Challenge string `redis:"challenge"`
	ClientId  uint64 `redis:"client_id"`
}

type OAuthCode struct {
	UserId uint64 `redis:"id"`
}

type OAuthTokens struct {
	Access        string
	Refresh       string
	ClientAccess  string
	ClientRefresh string
}

type FieldOptions struct {
	HasVerified bool
}

func (o *OAuthSession) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OAuthSession) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

func (o *OAuthCode) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OAuthCode) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
