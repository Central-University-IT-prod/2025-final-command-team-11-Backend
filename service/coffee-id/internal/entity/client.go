package entity

import (
	"encoding/json"

	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	"github.com/nikitaSstepanov/tools/client/pg"
)

type Client struct {
	Id          uint64        `redis:"id"`
	Name        string        `redis:"name"`
	ClientId    string        `redis:"client_id"`
	Secret      string        `redis:"secret"`
	RedirectUri string        `redis:"redirect_uri"`
	Fields      []types.Field `redis:"fields"`
	UserId      uint64        `redis:"user_id"`
}

func (c *Client) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Client) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *Client) Scan(r pg.Row) error {
	return r.Scan(
		&c.Id,
		&c.Name,
		&c.ClientId,
		&c.Secret,
		&c.RedirectUri,
		&c.Fields,
		&c.UserId,
	)
}
