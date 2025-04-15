package verification

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type VerificationUseCase interface {
	CheckVerify(c ctx.Context, id string) (*entity.Verification, e.Error)
	Verify(c ctx.Context, id string, passport *entity.Image) e.Error
}
