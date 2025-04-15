package coffeeid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/models"
)

type UsersRepo struct {
	coffeeIdBaseUrl string
}

func NewUserRepo(
	coffeeIdBaseUrl string,
) *UsersRepo {
	return &UsersRepo{
		coffeeIdBaseUrl: coffeeIdBaseUrl,
	}
}

func (ur *UsersRepo) GetById(ctx context.Context, id uuid.UUID) (models.User, error) {
	op := "coffee-id.UserRepo.GetById"

	url := fmt.Sprintf("%s/account/%s", ur.coffeeIdBaseUrl, id.String())

	resp, err := http.Get(url)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: make request: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return models.User{}, models.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("%s: get user by id: unexpected code %d", op, resp.StatusCode)
	}

	var user models.User
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: json.Decode: %w", op, err)
	}

	return user, nil
}
