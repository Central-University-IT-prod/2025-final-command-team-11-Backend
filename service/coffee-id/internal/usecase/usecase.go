package usecase

import (
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/admin"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/pkg/account"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/pkg/auth"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/storage"
	"github.com/nikitaSstepanov/tools"
	"github.com/nikitaSstepanov/tools/httper"
)

type UseCase struct {
	Account *account.Account
	Auth    *auth.Auth
}

type Config struct {
	Jwt   auth.JwtOptions  `yaml:"jwt"`
	Admin httper.ClientCfg `yaml:"admin"`
}

func New(storage *storage.Storage, cfg *Config) *UseCase {
	jwtAuth := auth.NewJwt(&cfg.Jwt)
	adm := admin.New(&cfg.Admin)
	coder := tools.Coder()

	account := account.New(
		&account.Storages{
			User: storage.Users,
			Code: storage.Codes,
		},
		&account.UseCases{
			Admin: adm,
			Jwt:   jwtAuth,
			Coder: coder,
		},
	)

	auth := auth.New(
		&auth.Storages{
			User: storage.Users,
		},
		&auth.UseCases{
			Jwt:   jwtAuth,
			Coder: coder,
		},
	)

	return &UseCase{
		Account: account,
		Auth:    auth,
	}
}
