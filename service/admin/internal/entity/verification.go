package entity

import "github.com/nikitaSstepanov/tools/client/pg"

type Verification struct {
	UserId        string
	PassportImage string
}

type Image struct {
	Name        string
	Buffer      []byte
	Size        int64
	ContentType string
}

func (v *Verification) Scan(r pg.Row) error {
	return r.Scan(
		&v.UserId,
		&v.PassportImage,
	)
}
