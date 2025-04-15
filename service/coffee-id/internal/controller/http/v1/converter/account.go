package converter

import (
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/dto"
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
)

func DtoUser(user *entity.User) *dto.Account {
	return &dto.Account{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}

func DtoAnswer(user *entity.User, token string) *dto.AccountAnswer {
	return &dto.AccountAnswer{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}
}

func DtoUserForOAuth(user *entity.User, opts *entity.FieldOptions) map[string]interface{} {
	data := make(map[string]interface{})

	if user.Id != "" {
		data["id"] = user.Id
	}

	if user.Email != "" {
		data["email"] = user.Email
	}

	if user.Name != "" {
		data["name"] = user.Name
	}

	if len(user.Role) != 0 {
		data["role"] = user.Role
	}

	if opts.HasVerified {
		data["verified"] = user.Verified
	}

	return data
}

func EntityCreate(create dto.CreateUser) *entity.User {
	return &entity.User{
		Email:    create.Email,
		Name:     create.Name,
		Password: create.Password,
	}
}

func EntityUpdate(update dto.UpdateUser) *entity.User {
	return &entity.User{
		Email:    update.Email,
		Name:     update.Name,
		Password: update.Password,
	}
}

func EntitySetRole(role dto.SetRole) *entity.User {
	return &entity.User{
		Id:   role.Id,
		Role: types.Role(role.Role),
	}
}

func EntityDelete(delete dto.DeleteUser) *entity.User {
	return &entity.User{
		Password: delete.Password,
	}
}
