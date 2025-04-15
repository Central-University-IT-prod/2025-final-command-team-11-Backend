package entity

import (
	"encoding/json"
	"time"

	"github.com/nikitaSstepanov/tools/client/pg"
)

type Yandex struct {
	YandexId string    `redis:"yandex_id"`
	Email    string    `redis:"email"`
	Name     string    `redis:"name"`
	Birthday time.Time `redis:"birthday"`
	UserId   uint64    `redis:"user_id"`
}

func (y *Yandex) MarshalBinary() ([]byte, error) {
	return json.Marshal(y)
}

func (y *Yandex) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, y)
}

func (y *Yandex) Scan(r pg.Row) error {
	return r.Scan(
		&y.YandexId,
		&y.Email,
		&y.Name,
		&y.Birthday,
		&y.UserId,
	)
}
