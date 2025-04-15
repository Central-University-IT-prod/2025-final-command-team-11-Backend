package guest

import (
	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	conv "REDACTED/team-11/backend/admin/internal/controller/http/v1/converter"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/dto"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/validator"
	resp "REDACTED/team-11/backend/admin/internal/controller/response"
	ct "REDACTED/team-11/backend/admin/pkg/utils/controller"
)

type Guest struct {
	usecase GuestUseCase
}

func New(uc GuestUseCase) *Guest {
	return &Guest{
		usecase: uc,
	}
}

// @Summary Create invite
// @Description Create invite to room.
// @Tags Guests
// @Security Bearer
// @Param id path string true "Booking id" Format(uuid)
// @Param Email body dto.GuestId true	"User Email"
// @Success 204 "Successfull invite"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Not your booking"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/booking/{id}/guests [post]
func (g *Guest) Create(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.GetString("userId")

	bookId := c.Param("id")

	if err := validator.UUID(bookId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	var body dto.GuestId

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, e.BadInputErr.WithErr(err))
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	err := g.usecase.Create(ctx, bookId, body.UserId, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}

// @Summary Get invites
// @Description Get invites of room.
// @Tags Guests
// @Security Bearer
// @Param id path string true "Booking id" Format(uuid)
// @Success 200 {object} []dto.Guest "Successfull get"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Not your booking"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/booking/{id}/guests [get]
func (g *Guest) Get(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.GetString("userId")

	bookId := c.Param("id")

	if err := validator.UUID(bookId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	guests, err := g.usecase.Get(ctx, bookId, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := make([]*dto.Guest, 0)

	for _, guest := range guests {
		result = append(result, conv.DtoGuest(guest))
	}

	c.JSON(httper.StatusOK, result)
}

// @Summary Delete invite
// @Description Delete invite to room.
// @Tags Guests
// @Security Bearer
// @Param id path string true "Booking id" Format(uuid)
// @Param email path string true "Email" Format(email)
// @Success 200 {object} []dto.Guest "Successfull get"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Not your booking"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/booking/{id}/guests/{email} [delete]
func (g *Guest) Delete(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.GetString("userId")

	bookId := c.Param("id")

	if err := validator.UUID(bookId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	userId := c.Param("userId")

	if err := validator.Email(userId); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	err := g.usecase.Delete(ctx, bookId, userId, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}
