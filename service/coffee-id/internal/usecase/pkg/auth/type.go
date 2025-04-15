package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

type JwtOptions struct {
	Audience   []string `yaml:"audience" env:"JWT_AUDIENCE"`
	Issuer     string   `yaml:"issuer"   env:"JWT_ISSUER"`
	AccessKey  string   `env:"JWT_ACCESS_SECRET"`
	RefreshKey string   `env:"JWT_REFRESH_SECRET"`
}

type Claims struct {
	Id   string     `json:"id"`
	Role types.Role `json:"role"`
	jwt.RegisteredClaims
}

type UseCases struct {
	YandexOAuth YandexOAuth
	Jwt         *Jwt
	Coder       *coder.Coder
}

type Storages struct {
	User   UserStorage
	OAuth  OAuthStorage
}

type UserStorage interface {
	GetById(ctx ctx.Context, id string) (*entity.User, e.Error)
	GetByEmail(ctx ctx.Context, email string) (*entity.User, e.Error)
	Create(ctx ctx.Context, user *entity.User) e.Error
	Verify(ctx ctx.Context, user *entity.User) e.Error
}

type YandexOAuth interface {
	GetUser(c ctx.Context, code string) (*entity.Yandex, e.Error)
}

type OAuthStorage interface {
	SetSession(c ctx.Context, key string, session *entity.OAuthSession) e.Error
	GetSession(c ctx.Context, key string) (*entity.OAuthSession, e.Error)
	SetCode(c ctx.Context, key string, code *entity.OAuthCode) e.Error
	GetCode(c ctx.Context, key string) (*entity.OAuthCode, e.Error)
}
