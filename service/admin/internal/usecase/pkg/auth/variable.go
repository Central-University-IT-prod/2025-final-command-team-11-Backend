package auth

import (
	"github.com/golang-jwt/jwt/v5"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

type JwtOptions struct {
	Audience  []string `yaml:"audience" env:"JWT_AUDIENCE"`
	Issuer    string   `yaml:"issuer"   env:"JWT_ISSUER"`
	AccessKey string   `env:"JWT_SECRET"`
}

type Claims struct {
	Id   string     `json:"id"`
	Role types.Role `json:"role"`
	jwt.RegisteredClaims
}
