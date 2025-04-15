package verification

import (
	"strings"

	"github.com/google/uuid"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Verification struct {
	verification VerificationStorage
	id           IdUseCase
}

func New(verfy VerificationStorage, id IdUseCase) *Verification {
	return &Verification{
		verification: verfy,
		id:           id,
	}
}

func (v *Verification) CheckVerify(c ctx.Context, id string) (*entity.Verification, e.Error) {
	_, err := v.id.GetUserById(c, id)
	if err != nil {
		return nil, err
	}

	return v.verification.CheckVerify(c, id)
}

func (v *Verification) Verify(c ctx.Context, id string, passport *entity.Image) e.Error {
	_, err := v.id.GetUserById(c, id)
	if err != nil {
		return err
	}

	data, err := v.verification.Get(c, id)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	if err == nil {
		image := &entity.Image{
			Name:        data.PassportImage,
			Buffer:      passport.Buffer,
			Size:        passport.Size,
			ContentType: passport.ContentType,
		}

		return v.verification.UpdateData(c, image)
	}

	name := passport.Name
	parts := strings.Split(name, ".")

	if len(parts) < 2 {
		return e.New("Bad file name.", e.BadInput)
	}

	imageName := uuid.NewString()
	parts[0] = imageName

	name = strings.Join(parts, ".")
	passport.Name = name

	return v.verification.Verify(c, id, passport)
}
