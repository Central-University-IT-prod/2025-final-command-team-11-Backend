package order

import (
	"strconv"

	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	conv "REDACTED/team-11/backend/admin/internal/controller/http/v1/converter"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/dto"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/validator"
	resp "REDACTED/team-11/backend/admin/internal/controller/response"
	ct "REDACTED/team-11/backend/admin/pkg/utils/controller"
)

type Order struct {
	usecase OrderUseCase
}

func New(uc OrderUseCase) *Order {
	return &Order{
		usecase: uc,
	}
}

// @Summary Set order completed
// @Description Set order. Only for ADMINs
// @Tags Orders
// @Security Bearer
// @Param id    path     string  true  "user id"  Format(uuid)
// @Success 204 "Successful set order status"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/orders/{id} [post]
func (o *Order) SetStatus(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	err := o.usecase.ChangeStatus(ctx, id, true)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}

// @Summary Get orders
// @Description Get orders with pagination and filters. Only for ADMINs.
// @Tags Orders
// @Security Bearer
// @Param page      query int  false  "Page"
// @Param size      query int  false  "Size"
// @Param completed query string  false  "Completed. Complteted must be 'true', 'false' or ''"
// @Success 200 {object} dto.Orders "Succesfull get of orders"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/orders [get]
func (o *Order) Get(c *gin.Context) {
	ctx := ct.GetCtx(c)

	page, parseErr := strconv.ParseInt(c.DefaultQuery("page", "0"), 10, 64)
	if parseErr != nil || page < 0 {
		resp.AbortErrMsg(c, e.New("Page must be integer", e.BadInput))
		return
	}

	size, parseErr := strconv.ParseInt(c.DefaultQuery("size", "5"), 10, 64)
	if parseErr != nil || size < 0 {
		resp.AbortErrMsg(c, e.New("Page must be integer", e.BadInput))
		return
	}

	filter := c.DefaultQuery("completed", "")
	if filter != "true" && filter != "false" && filter != "" {
		resp.AbortErrMsg(c, e.New(`Complteted must be "true", "false" or ""`, e.BadInput))
		return
	}

	orders, count, err := o.usecase.GetOrders(ctx, int(page), int(size), filter)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := make([]*dto.Order, 0)

	for _, order := range orders {
		result = append(result, conv.DtoOrder(order))
	}

	c.JSON(httper.StatusOK, dto.Orders{Values: result, Count: count})
}

// @Summary Get stats
// @Description Get stats of order creations. Only for ADMINs.
// @Tags Orders
// @Param filter query string  false  "Filter"
// @Success 200 {object} dto.Stats
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/orders/stats [get]
func (b *Order) Stats(c *gin.Context) {
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
