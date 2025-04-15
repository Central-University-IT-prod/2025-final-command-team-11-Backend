package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/pkg/logger"
)

var (
	bookingsTable = "booking"
)

type BookingsRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewBookingsRepo(db *sqlx.DB) *BookingsRepo {
	return &BookingsRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (br *BookingsRepo) Create(ctx context.Context, input dto.BookingCreateDto) (models.Booking, error) {
	op := "postgres.BookingsRepo.Create"

	query, args, err := br.sq.
		Insert(bookingsTable).
		Columns("entity_id", "user_id", "time_from", "time_to").
		Values(input.EntityId, input.UserId, input.TimeFrom, input.TimeTo).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res models.Booking
	if err := br.db.GetContext(ctx, &res, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return models.Booking{}, models.ErrBookingEntityNotFound
			}
		}

		return models.Booking{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return res, nil
}

func (br *BookingsRepo) GetById(ctx context.Context, id uuid.UUID) (models.Booking, error) {
	op := "postgres.BookingsRepo.GetById"

	query, args, err := br.sq.
		Select("*").
		From(bookingsTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res models.Booking
	if err := br.db.GetContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Booking{}, models.ErrBookingNotFound
		}

		return models.Booking{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return res, nil
}

func (br *BookingsRepo) ListForUser(ctx context.Context, userId uuid.UUID) ([]models.Booking, error) {
	op := "postgres.BookingsRepo.ListForUser"

	query, args, err := br.sq.
		Select("*").
		From(bookingsTable).
		Where(sq.Eq{
			"user_id": userId,
		}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res []models.Booking
	if err := br.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return res, nil
}

func (br *BookingsRepo) ListAll(ctx context.Context) ([]models.Booking, error) {
	op := "postgres.BookingsRepo.ListAll"

	query, args, err := br.sq.
		Select("*").
		From(bookingsTable).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res []models.Booking
	if err := br.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return res, nil
}

func (br *BookingsRepo) Update(ctx context.Context, input dto.BookingUpdateDto) (models.Booking, error) {
	op := "postgres.BookingsRepo.Update"

	qb := br.sq.Update(bookingsTable)

	var fieldsUpdates int

	if input.TimeFrom != nil {
		fieldsUpdates++
		qb = qb.Set("time_from", *input.TimeFrom)
	}
	if input.TimeTo != nil {
		fieldsUpdates++
		qb = qb.Set("time_to", *input.TimeTo)
	}

	if fieldsUpdates == 0 {
		qb = qb.Set("updated_at", time.Now().UTC())
	}

	query, args, err := qb.
		Where(sq.Eq{"id": input.BookingId}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res models.Booking
	if err := br.db.GetContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Booking{}, models.ErrBookingNotFound
		}

		return models.Booking{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return res, nil
}

func (br *BookingsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "postgres.BookingsRepo.Delete"

	query, args, err := br.sq.
		Delete(bookingsTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	res, err := br.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: res.RowsAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return models.ErrBookingNotFound
	}

	return nil
}

func (br *BookingsRepo) ListIntersectedForUser(ctx context.Context, userId uuid.UUID, timeFrom, timeTo time.Time) ([]models.Booking, error) {
	op := "postgres.BookingsRepo.ListInterSectedForUser"

	query, args, err := br.sq.
		Select("*").
		From(bookingsTable).
		Where(sq.And{
			sq.Eq{"user_id": userId},
			sq.Or{
				sq.And{
					sq.Gt{"time_from": timeFrom},
					sq.Lt{"time_from": timeTo},
				},
				sq.And{
					sq.Gt{"time_to": timeFrom},
					sq.Lt{"time_to": timeTo},
				},
				sq.Eq{
					"time_to":   timeTo,
					"time_from": timeFrom,
				},
				sq.And{
					sq.Lt{"time_from": timeFrom},
					sq.Gt{"time_to": timeTo},
				},
			},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	logger.FromCtx(ctx).Debug("query: " + query)

	var res []models.Booking
	if err := br.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return res, nil
}

func (br *BookingsRepo) ListIntersectedForEntity(ctx context.Context, entityId uuid.UUID, timeFrom, timeTo time.Time) ([]models.Booking, error) {
	op := "postgres.BookingsRepo.ListInterSected"

	query, args, err := br.sq.
		Select("*").
		From(bookingsTable).
		Where(sq.And{
			sq.Eq{"entity_id": entityId},
			sq.Or{
				sq.And{
					sq.Gt{"time_from": timeFrom},
					sq.Lt{"time_from": timeTo},
				},
				sq.And{
					sq.Gt{"time_to": timeFrom},
					sq.Lt{"time_to": timeTo},
				},
				sq.Eq{
					"time_to":   timeTo,
					"time_from": timeFrom,
				},
				sq.And{
					sq.Lt{"time_from": timeFrom},
					sq.Gt{"time_to": timeTo},
				},
			},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res []models.Booking
	if err = br.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return res, nil
}
