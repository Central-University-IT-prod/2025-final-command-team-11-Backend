package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
)

var (
	ordersTable = "orders"
)

type OrdersRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (or *OrdersRepo) Create(ctx context.Context, input dto.OrderCreateDto) (models.Order, error) {
	op := "postgres.OrdersRepo.Create"

	query, args, err := or.sq.
		Insert(ordersTable).
		Columns("booking_id", "thing").
		Values(input.BookingId, input.Thing).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var creared models.Order
	if err := or.db.GetContext(ctx, &creared, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return models.Order{}, models.ErrBookingNotFound
			}
		}

		return models.Order{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return creared, nil
}

func (or *OrdersRepo) GetById(ctx context.Context, id uuid.UUID) (models.Order, error) {
	op := "postgres.OrdersRepo.GetById"

	query, args, err := or.sq.
		Select("*").
		From(ordersTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var order models.Order
	if err := or.db.GetContext(ctx, &order, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, models.ErrOrderNotFound
		}

		return models.Order{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return order, nil
}

func (or *OrdersRepo) GetForBooking(ctx context.Context, bookingId uuid.UUID) ([]models.Order, error) {
	op := "postgres.OrdersRepo.GetForBooking"

	query, args, err := or.sq.
		Select("*").
		From(ordersTable).
		Where(sq.Eq{"booking_id": bookingId}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res []models.Order
	if err := or.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return res, nil
}

func (or *OrdersRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "postgres.OrdersRepo.Delete"

	query, args, err := or.sq.
		Delete(ordersTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	res, err := or.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	if rowsAffected == 0 {
		return models.ErrOrderNotFound
	}

	return nil
}
