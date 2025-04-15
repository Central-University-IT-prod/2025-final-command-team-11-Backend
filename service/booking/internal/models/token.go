package models

import "github.com/google/uuid"

type Token struct {
	UserId uuid.UUID
	Role   Role
}

func TokenFromCliams(claims map[string]any) (Token, error) {
	userIdStr, ok := claims["id"].(string)
	if !ok {
		return Token{}, ErrInvalidToken
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return Token{}, ErrInvalidToken
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return Token{}, ErrInvalidToken
	}

	role, err := ParseRole(roleStr)
	if err != nil {
		return Token{}, ErrInvalidToken
	}

	return Token{
		UserId: userId,
		Role:   role,
	}, nil
}
