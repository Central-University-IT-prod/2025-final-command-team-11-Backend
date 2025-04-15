package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"REDACTED/team-11/backend/booking/internal/models"
)

var (
	bookingEntitiesTable = "booking_entity"
)

type BookingEntitiesRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewBookingEntitiesRepo(db *sqlx.DB) *BookingEntitiesRepo {
	return &BookingEntitiesRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (ber *BookingEntitiesRepo) GetById(ctx context.Context, id uuid.UUID) (models.BookingEntity, error) {
	op := "postgres.BookingEntitiesRepo.GetById"

	query, args, err := ber.sq.
		Select("*").
		From(bookingEntitiesTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.BookingEntity{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res models.BookingEntity
	if err = ber.db.GetContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BookingEntity{}, models.ErrBookingEntityNotFound
		}

		return models.BookingEntity{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return res, nil
}

func (ber *BookingEntitiesRepo) GetForFloor(ctx context.Context, floorId uuid.UUID) ([]models.BookingEntity, error) {
	op := "postgres.BookingEntitiesRepo.GetForFloor"

	query, args, err := ber.sq.
		Select("*").
		From(bookingEntitiesTable).
		Where(sq.Eq{"floor_id": floorId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res []models.BookingEntity
	if err := ber.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return res, nil
}
