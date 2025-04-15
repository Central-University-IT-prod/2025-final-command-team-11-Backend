package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/internal/transport/http/v1/security"
	"REDACTED/team-11/backend/booking/pkg/logger"
	api "REDACTED/team-11/backend/booking/pkg/ogen"
	"go.uber.org/zap"
)

type WorkloadsUsecase interface {
	GetForFloor(ctx context.Context, floorId uuid.UUID, timeFrom, timeTo time.Time, userId uuid.UUID) (models.FloorWorkload, error)
	Get(ctx context.Context, entityId uuid.UUID, timeFrom, timeTo time.Time) (models.Workload, error)
}

type WorkloadsHandler struct {
	usecase WorkloadsUsecase
}

func NewWorkloadsHandler(usecase WorkloadsUsecase) *WorkloadsHandler {
	return &WorkloadsHandler{
		usecase: usecase,
	}
}

// GetWorkload implements getWorkload operation.
//
// Get workload booking entity.
//
// GET /workloads/{entityId}
func (wh *WorkloadsHandler) GetWorkload(ctx context.Context, params api.GetWorkloadParams) (api.GetWorkloadRes, error) {

	if params.TimeFrom >= params.TimeTo {
		return &api.Response400{
			Message: api.NewOptString("time_from must be before time_to"),
		}, nil
	}

	if int64(params.TimeFrom)%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_from must be multiple of 15 minutes"),
		}, nil
	}

	if int64(params.TimeTo)%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_to must be multiple of 15 minutes"),
		}, nil
	}

	timeFrom := time.Unix(int64(params.TimeFrom), 0).UTC()
	timeTo := time.Unix(int64(params.TimeTo), 0).UTC()

	workload, err := wh.usecase.Get(ctx, params.EntityId, timeFrom, timeTo)
	if err != nil {
		if errors.Is(err, models.ErrBookingEntityNotFound) {
			return &api.Response404{}, nil
		}

		logger.FromCtx(ctx).Error("get workload", zap.Error(err))
		return nil, err
	}

	res := api.Workload(make([]api.WorkloadItem, 0, len(workload)))
	for _, snapshot := range workload {
		res = append(res, api.WorkloadItem{
			Time:   api.Time(snapshot.Time.Unix()),
			IsFree: snapshot.IsFree,
		})
	}

	return &res, nil
}

// GetFloorWorkload implements getFloorWorkload operation.
//
// Возвращает информацию о нагрузке на указанный этаж
// за указанный период времени.
//
// GET /workloads/floors/{floorId}
func (wh *WorkloadsHandler) GetFloorWorkload(ctx context.Context, params api.GetFloorWorkloadParams) (api.GetFloorWorkloadRes, error) {
	token := security.TokenFromCtx(ctx)

	if params.TimeFrom >= params.TimeTo {
		return &api.Response400{
			Message: api.NewOptString("time_from must be before time_to"),
		}, nil
	}

	if int64(params.TimeFrom)%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_from must be multiple of 15 minutes"),
		}, nil
	}

	if int64(params.TimeTo)%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_to must be multiple of 15 minutes"),
		}, nil
	}

	timeFrom := time.Unix(int64(params.TimeFrom), 0).UTC()
	timeTo := time.Unix(int64(params.TimeTo), 0).UTC()

	floorWorkload, err := wh.usecase.GetForFloor(ctx, params.FloorId, timeFrom, timeTo, token.UserId)

	if err != nil {
		if errors.Is(err, models.ErrFloorNotFound) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceFloor),
			}, nil
		}
		logger.FromCtx(ctx).Error("get floor workload", zap.Error(err))
		return nil, err
	}

	res := make(api.FloorWorkload, 0, len(floorWorkload))
	for _, entityWorkload := range floorWorkload {
		res = append(res, api.FloorWorkloadItem{
			Entity: convertBookingEntity(entityWorkload.Entity),
			IsFree: entityWorkload.IsFree,
		})
	}

	return &res, nil
}
