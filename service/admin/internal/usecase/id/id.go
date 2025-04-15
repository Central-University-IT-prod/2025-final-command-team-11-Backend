package id

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Id struct {
	client *httper.Client
}

func New(cfg *httper.ClientCfg) *Id {
	return &Id{
		client: httper.NewClient(cfg),
	}
}

func (i *Id) GetUser(c ctx.Context, email string) (*entity.User, e.Error) {
	var user entity.User

	req, err := httper.NewReq(&httper.Params{
		Method:        httper.GetMethod,
		Url:           "/account/email/" + email,
		Unmarshal:     true,
		UnmarshalTo:   &user,
		UnmarshalType: httper.JsonType,
	})
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	response, err := i.client.Do(req)
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	if response.StatusCode == 404 {
		return nil, e.New("User wasn`t found", e.NotFound)
	} else if response.StatusCode != 200 {
		return nil, e.InternalErr
	}

	return &user, nil
}

func (i *Id) GetUserById(c ctx.Context, id string) (*entity.User, e.Error) {
	var user entity.User

	req, err := httper.NewReq(&httper.Params{
		Method:        httper.GetMethod,
		Url:           "/account/" + id,
		Unmarshal:     true,
		UnmarshalTo:   &user,
		UnmarshalType: httper.JsonType,
	})
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	response, err := i.client.Do(req)
	if err != nil {
		return nil, e.InternalErr.WithErr(err)
	}

	if response.StatusCode == 404 {
		return nil, e.New("User wasn`t found", e.NotFound).WithTag("coffee", true)
	} else if response.StatusCode != 200 {
		return nil, e.InternalErr
	}

	return &user, nil
}
