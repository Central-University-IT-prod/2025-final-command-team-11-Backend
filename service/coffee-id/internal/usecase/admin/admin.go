package admin

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
)

type Admin struct {
	client *httper.Client
}

func New(cfg *httper.ClientCfg) *Admin {
	return &Admin{
		client: httper.NewClient(cfg),
	}
}

type User struct {
	Verified bool `json:"verified"`
}

func (i *Admin) GetUser(c ctx.Context, email string, token string) (*User, e.Error) {
	var user User

	req, err := httper.NewReq(&httper.Params{
		Method:        httper.GetMethod,
		Url:           "/verification/" + email + "/check",
		Unmarshal:     true,
		UnmarshalTo:   &user,
		UnmarshalType: httper.JsonType,
	})
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	req.Header.Add("Authorization", token)

	_, err = i.client.Do(req)
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	return &user, nil
}
