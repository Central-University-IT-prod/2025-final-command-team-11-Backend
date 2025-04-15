package booking_entity

import (
	"time"

	"github.com/gin-gonic/gin"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
	conv "REDACTED/team-11/backend/admin/internal/controller/http/v1/converter"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/dto"
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/validator"
	resp "REDACTED/team-11/backend/admin/internal/controller/response"
	"REDACTED/team-11/backend/admin/internal/entity"
	ct "REDACTED/team-11/backend/admin/pkg/utils/controller"
)

type BookingEntity struct {
	usecase EntityuseCase
}

func New(uc EntityuseCase) *BookingEntity {
	return &BookingEntity{
		usecase: uc,
	}
}

// @Summary Get floors
// @Description Get list of floor
// @Tags Entity
// @Success 200 {object} []dto.FloorEntity "Succesful get of floors"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/layout/floors [get]
func (b *BookingEntity) GetFloors(c *gin.Context) {
	ctx := ct.GetCtx(c)

	floors, err := b.usecase.GetFloors(ctx)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := make([]*dto.FloorEntity, 0)

	for _, floor := range floors {
		result = append(result, conv.DtoFloor(floor))
	}

	c.JSON(httper.StatusOK, result)
}

// @Summary Get entities for floor
// @Description Get entities for floor
// @Tags Entity
// @Param id path string true "Floor id" Format(uuid)
// @Success 200 {object} []dto.BookingEntity "Successful get entities"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/layout/floors/{id} [get]
func (b *BookingEntity) GetEntities(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	bookings, err := b.usecase.GetEntities(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := make([]*dto.BookingEntity, 0)

	for _, booking := range bookings {
		result = append(result, conv.DtoEntity(booking))
	}

	c.JSON(httper.StatusOK, result)
}

// @Summary Save layout
// @Description Save layout. Only for ADMINs
// @Tags Entity
// @Accept json
// @Param upsert body dto.UpsertFloor true	"Upsert data"
// @Success 200 {object} resp.Message "OK"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/layout/floors [post]
func (b *BookingEntity) Save(c *gin.Context) {
	ctx := ct.GetCtx(c)

	var body dto.UpsertFloor

	if err := c.ShouldBindJSON(&body); err != nil {
		resp.AbortErrMsg(c, e.BadInputErr.WithErr(err))
		return
	}

	if err := validator.Struct(body); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	for _, booking := range body.Entities {
		if err := validator.Struct(booking, validator.Booking); err != nil {
			resp.AbortErrMsg(c, err)
			return
		}
	}

	bookings := make([]*entity.BookingEntity, 0)
	curTime := time.Now().UTC()

	for _, booking := range body.Entities {
		toAdd := &entity.BookingEntity{
			Id:        booking.Id,
			Type:      booking.Type,
			Title:     booking.Title,
			X:         booking.X,
			Y:         booking.Y,
			FloorId:   booking.FloorId,
			Width:     booking.Width,
			Height:    booking.Height,
			Capacity:  booking.Capacity,
			CreatedAt: curTime,
			UpdatedAt: curTime,
		}

		bookings = append(bookings, toAdd)
	}

	floor := &entity.FloorEntity{
		Id:        body.Id,
		Name:      body.Name,
		CreatedAt: curTime,
		UpdatedAt: curTime,
	}

	err := b.usecase.Save(ctx, bookings, floor)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusOK, resp.NewMessage("OK"))
}

// @Summary Delete floor
// @Description Delete floor. Only for ADMINs
// @Tags Entity
// @Param id path string true "Floor id" Format(uuid)
// @Success 204 "Successful delete"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/layout/floors/{id} [delete]
func (b *BookingEntity) DeleteFloor(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	if err := b.usecase.DeleteFloor(ctx, id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(httper.StatusNoContent, nil)
}

// @Summary Get entity by id
// @Description Get entity by id
// @Tags Entity
// @Param id path string true "Entity id"  Format(uuid)
// @Success 200 {object} dto.BookingEntity "ok"
// @Failure 401 {object} resp.JsonError "Unauth"
// @Failure 403 {object} resp.JsonError "Invalid role"
// @Failure 400 {object} resp.JsonError "Id must be uuid"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /admin/layout/entities/{id} [get]
func (b *BookingEntity) EntityById(c *gin.Context) {
	ctx := ct.GetCtx(c)

	id := c.Param("id")

	if err := validator.UUID(id); err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	booking, err := b.usecase.GetEntity(ctx, id)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.DtoEntity(booking)

	c.JSON(httper.StatusOK, result)
}
