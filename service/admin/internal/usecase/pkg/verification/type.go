package verification

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type VerificationStorage interface {
	CheckVerify(c ctx.Context, id string) (*entity.Verification, e.Error)
	Verify(c ctx.Context, id string, image *entity.Image) e.Error
	UpdateData(c ctx.Context, image *entity.Image) e.Error
	Get(c ctx.Context, id string) (*entity.Verification, e.Error)
}

type IdUseCase interface {
	GetUser(c ctx.Context, email string) (*entity.User, e.Error)
	GetUserById(c ctx.Context, id string) (*entity.User, e.Error)
}
