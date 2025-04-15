package booking

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	conv "REDACTED/team-11/backend/admin/internal/controller/http/v1/converter"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/validator"
	resp "REDACTED/team-11/backend/admin/internal/controller/response"
	ct "REDACTED/team-11/backend/admin/pkg/utils/controller"
)

type Booking struct {
	usecase BookingUseCase
}

func New(uc BookingUseCase) *Booking {
	return &Booking{
		usecase: uc,
	}
}

// @Summary Check user access
// @Description Check status of user booking or invitation for nearest 12 hours. Avaliable only for ADMINs
// @Tags Booking
// @Security Bearer
// @Param id path string true  "User id"  Format(uuid)
// @Success 200 {object} dto.BookingAccess "Successful check access. Status can has value READY, NOT_READY or PENDING"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/booking/{id}/access [get]
func (b *Booking) CheckAccess(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	booking, err := b.usecase.CheckAccess(ctx, id)
	if err != nil && err.GetCode() != e.NotFound {
		resp.AbortErrMsg(c, err)
		return
	}

	if err != nil {
		result := conv.DtoAccess("NOT_READY", nil)

		c.JSON(httper.StatusOK, result)
		return
	}
	fmt.Println(booking)
	if booking.TimeFrom.After(time.Now().UTC()) {
		result := conv.DtoAccess("PENDING", booking)

		c.JSON(httper.StatusOK, result)
		return
	}

	result := conv.DtoAccess("READY", booking)

	c.JSON(httper.StatusOK, result)
}

// @Summary Get stats
// @Description Get stats for booking creations. Avaliable only for ADMINs
// @Tags Booking
// @Security Bearer
// @Param filter query string false "Parametr for stats specify. Must be 'day', 'week' or 'month'"
// @Success 200 {object} dto.Stats "Successful get of stats"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/booking/stats [get]
func (b *Booking) Stats(c *gin.Context) {
	ctx := ct.GetCtx(c)

	filter := c.DefaultQuery("filter", "day")
	if filter != "day" && filter != "month" && filter != "week" {
		resp.AbortErrMsg(c, e.New(`Filter must be "day", "month" or "week"`, e.BadInput))
		return
	}

	cnt, err := b.usecase.Stats(ctx, filter)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusOK, gin.H{"count": cnt})
}
