package account

import (
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/dto"
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

type Account struct {
	user  UserStorage
	code  CodeStorage
	jwt   JwtUseCase
	adm   Admin
	coder *coder.Coder
}

func New(store *Storages, uc *UseCases) *Account {
	return &Account{
		user:  store.User,
		code:  store.Code,
		jwt:   uc.Jwt,
		coder: uc.Coder,
		adm:   uc.Admin,
	}
}

func (a *Account) GetList(c ctx.Context, page, size int, token string) ([]*dto.AP, int, e.Error) {
	o, err := a.user.Getc(c)
	if err != nil {
		return nil, 0, err
	}

	res, err := a.user.Get(c, size, page*size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.AP, 0)

	for _, acc := range res {
		dd, err := a.adm.GetUser(c, acc.Id, token)
		if err != nil {
			return nil, 0, err
		}

		to := &dto.AP{
			Id:       acc.Id,
			Email:    acc.Email,
			Name:     acc.Name,
			Verified: dd.Verified,
		}

		result = append(result, to)
	}

	return result, len(o), nil
}

func (a *Account) GetByEmail(c ctx.Context, email string) (*entity.User, e.Error) {
	return a.user.GetByEmail(c, email)
}

func (a *Account) Get(ctx ctx.Context, userId string) (*entity.User, e.Error) {
	user, err := a.user.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *Account) Create(ctx ctx.Context, user *entity.User) (*entity.Tokens, e.Error) {
	candidate, err := a.user.GetByEmail(ctx, user.Email)
	if err != nil && err.GetCode() != e.NotFound {
		return nil, err
	}

	if candidate != nil {
		return nil, conflictErr
	}

	hash, hashErr := a.coder.Hash(user.Password)
	if hashErr != nil {
		return nil, e.InternalErr.WithErr(err).WithCtx(ctx)
	}

	user.Password = hash

	err = a.user.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Role = defaultRole

	return a.getTokens(user)
}

func (a *Account) Update(ctx ctx.Context, user *entity.User, pass string) (*entity.User, e.Error) {
	if user.Email != "" {
		candidate, err := a.user.GetByEmail(ctx, user.Email)
		if err != nil && err.GetCode() != e.NotFound {
			return nil, err
		}

		if candidate != nil {
			return nil, conflictErr
		}

		user.Verified = false

		err = a.user.Verify(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	if user.Password != "" {
		old, err := a.user.GetById(ctx, user.Id)
		if err != nil {
			return nil, err
		}

		if err := a.coder.CompareHash(old.Password, pass); err != nil {
			return nil, badPassErr.WithErr(err)
		}

		hash, hashErr := a.coder.Hash(user.Password)
		if hashErr != nil {
			return nil, e.InternalErr.WithErr(err)
		}

		user.Password = hash
	}

	user, err := a.user.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *Account) AddRole(ctx ctx.Context, user *entity.User) e.Error {
	return a.user.AddRole(ctx, user)
}

func (a *Account) Delete(ctx ctx.Context, user *entity.User) e.Error {
	toDel, err := a.user.GetById(ctx, user.Id)
	if err != nil {
		return err
	}

	if err := a.coder.CompareHash(toDel.Password, user.Password); err != nil {
		return badPassErr.WithErr(err)
	}

	return a.user.Delete(ctx, user)
}

func (a *Account) getTokens(user *entity.User) (*entity.Tokens, e.Error) {
	var tokens entity.Tokens

	access, err := a.jwt.GenerateToken(user, false)
	if err != nil {
		return nil, err
	}

	refresh, err := a.jwt.GenerateToken(user, true)
	if err != nil {
		return nil, err
	}

	tokens.Access = access
	tokens.Refresh = refresh

	return &tokens, nil
}
