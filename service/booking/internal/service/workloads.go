package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/internal/repo"
	"REDACTED/team-11/backend/booking/pkg/logger"
	"go.uber.org/zap"
)

type WorkloadsService interface {
	GetForFloor(ctx context.Context, floorId uuid.UUID, timeFrom, timeTo time.Time, userId uuid.UUID) (models.FloorWorkload, error)
	Get(ctx context.Context, entityId uuid.UUID, timeFrom, timeTo time.Time) (models.Workload, error)
}

var (
	_ WorkloadsService = NewWorkloadService(nil, nil, nil)
)

var (
	intervalMinutes = 15
)

type workloadsServiceImpl struct {
	bookingEntitiesRepo repo.BookingEntitiesRepo
	bookingsRepo        repo.BookingsRepo
	floorsRepo          repo.FloorsRepo
}

func NewWorkloadService(
	bookingEntitiesRepo repo.BookingEntitiesRepo,
	bookingsRepo repo.BookingsRepo,
	floorsRepo repo.FloorsRepo,
) *workloadsServiceImpl {
	return &workloadsServiceImpl{
		bookingEntitiesRepo: bookingEntitiesRepo,
		bookingsRepo:        bookingsRepo,
		floorsRepo:          floorsRepo,
	}
}

func (ws *workloadsServiceImpl) Get(ctx context.Context, entityId uuid.UUID, timeFrom, timeTo time.Time) (models.Workload, error) {
	op := "service.workloadsServiceImpl.Get"

	bookingEntity, err := ws.bookingEntitiesRepo.GetById(ctx, entityId)
	if err != nil {
		if errors.Is(err, models.ErrBookingEntityNotFound) {
			return nil, models.ErrBookingEntityNotFound
		}

		return nil, fmt.Errorf("%s: bookingEntitiesRepo.GetById: %w", op, err)
	}

	intersected, err := ws.bookingsRepo.ListIntersectedForEntity(ctx, entityId, timeFrom, timeTo)
	if err != nil {
		return nil, fmt.Errorf("%s: bookingsRepo.ListInterSected: %w", op, err)
	}

	intersectedIds := make([]uuid.UUID, 0, len(intersected))
	for _, i := range intersected {
		intersectedIds = append(intersectedIds, i.Id)
	}

	logger.FromCtx(ctx).Debug("insersectedIds", zap.Any("", intersectedIds))

	snapshotsCount := (int(timeTo.Sub(timeFrom).Minutes()) / intervalMinutes) + 1
	takenPlaces := make([]int, snapshotsCount)

	for _, booking := range intersected {
		intersectionFrom := maxTime(booking.TimeFrom, timeFrom)
		intersectionTo := minTime(booking.TimeTo, timeTo)
		for snapshotTime := intersectionFrom; snapshotTime.Unix() <= intersectionTo.Unix(); snapshotTime = snapshotTime.Add(time.Duration(intervalMinutes) * time.Minute) {
			i := int(snapshotTime.Sub(timeFrom).Minutes()) / intervalMinutes
			takenPlaces[i] = takenPlaces[i] + 1
		}
	}

	res := make(models.Workload, 0, snapshotsCount)
	for i, takenPlacesSnapshot := range takenPlaces {
		var isFree bool
		switch bookingEntity.Type {
		case models.BookingEntityTypeOpenSpace:
			isFree = takenPlacesSnapshot < bookingEntity.Capacity
		case models.BookingEntityTypeRoom:
			isFree = takenPlacesSnapshot == 0
		}
		res = append(res, models.WorkloadItem{
			Time:   timeFrom.Add(time.Duration(i*intervalMinutes) * time.Minute),
			IsFree: isFree,
		})
	}

	return res, nil
}

func (ws *workloadsServiceImpl) GetForFloor(ctx context.Context, floorId uuid.UUID, timeFrom, timeTo time.Time, userId uuid.UUID) (models.FloorWorkload, error) {
	op := "service.workloadsServiceImpl.GetForFloor"

	_, err := ws.floorsRepo.GetById(ctx, floorId)
	if err != nil {
		if errors.Is(err, models.ErrFloorNotFound) {
			return nil, models.ErrFloorNotFound
		}

		return nil, fmt.Errorf("%s: floorsRepo.GetById: %w", op, err)
	}

	entities, err := ws.bookingEntitiesRepo.GetForFloor(ctx, floorId)
	if err != nil {
		return nil, fmt.Errorf("%s: bookingEntitiesRepo.GetForFloor: %w", op, err)
	}

	floorWorkload := make(models.FloorWorkload, 0, len(entities))
	for _, entity := range entities {
		workload, err := ws.Get(ctx, entity.Id, timeFrom, timeTo)
		if err != nil {
			return nil, fmt.Errorf("%s: ws.Get: %w", op, err)
		}
		isFree := true

		for _, snapshot := range workload {
			if !snapshot.IsFree {
				isFree = false
				break
			}
		}

		intersected, err := ws.bookingsRepo.ListIntersectedForUser(ctx, userId, timeFrom, timeTo)
		if err != nil {
			return nil, fmt.Errorf("%s: bookingsRepo.ListIntersectedForUser: %w", op, err)
		}

		if len(intersected) != 0 {
			isFree = false
		}

		floorWorkload = append(floorWorkload, models.FloorWorkloadItem{
			Entity: entity,
			IsFree: isFree,
		})
	}

	return floorWorkload, nil
}

func maxTime(times ...time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}

	res := times[0]
	for i := 1; i < len(times); i++ {
		if times[i].After(res) {
			res = times[i]
		}
	}

	return res
}

func minTime(times ...time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}

	res := times[0]
	for i := 1; i < len(times); i++ {
		if times[i].Before(res) {
			res = times[i]
		}
	}

	return res
}
