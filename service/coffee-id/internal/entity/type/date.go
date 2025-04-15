package types

import (
	"encoding/json"
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(bs []byte) error {
	var s string

	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation(time.DateOnly, s, time.UTC)
	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*d).Format(time.DateOnly))
}
