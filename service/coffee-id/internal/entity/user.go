package entity

import (
	"encoding/json"
	"time"

	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	"github.com/nikitaSstepanov/tools/client/pg"
)

type User struct {
	Id       string      `redis:"id"`
	Email    string      `redis:"email"`
	Name     string      `redis:"name"`
	Password string      `redis:"password"`
	Birthday time.Time   `redis:"birthday"`
	Role     types.Role  `redis:"roles"`
	Verified bool        `redis:"verified"`
	OAuth    types.OAuth `redis:"oauth"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *User) Scan(r pg.Row) error {
	var oauth *string

	err := r.Scan(
		&u.Id,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.Birthday,
		&u.Role,
		&u.Verified,
		&oauth,
	)
	if err != nil {
		return err
	}

	if oauth != nil {
		u.OAuth = types.OAuth(*oauth)
	}

	return nil
}
