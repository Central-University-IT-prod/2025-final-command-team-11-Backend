package order

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Order struct {
	postgres pg.Client
}

func New(postgres pg.Client) *Order {
	return &Order{
		postgres: postgres,
	}
}

func (o *Order) Update(c ctx.Context, order *entity.Order) e.Error {
	query, args, _ := sq.Update(orderTable).
		Set("completed", order.Completed).Set("updated_at", order.UpdatedAt).
		Where(sq.Eq{"id": order.Id}).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := o.postgres.Begin(c)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}
	defer tx.Rollback(c)

	if _, err := tx.Exec(c, query, args...); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	if err := tx.Commit(c); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (o *Order) Get(c ctx.Context) ([]*entity.OrderBooking, e.Error) {
	builder := sq.Select("*").From(fmt.Sprintf("%s AS o", orderTable)).
		Join(fmt.Sprintf("%s AS b ON o.booking_id = b.id", bookingTable)).
		Join(fmt.Sprintf("%s AS e ON b.entity_id = e.id", entityTable))

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()
	fmt.Println(query)
	rows, err := o.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Orders weren`t found.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	orders := make([]*entity.OrderBooking, 0)

	for rows.Next() {
		var order entity.Order
		var booking entity.Booking
		var bentity entity.BookingEntity

		err := rows.Scan(
			&order.Id,
			&order.BookingId,
			&order.Completed,
			&order.Thing,
			&order.CreatedAt,
			&order.UpdatedAt,
			&booking.Id,
			&booking.EntityId,
			&booking.UserId,
			&booking.TimeFrom,
			&booking.TimeTo,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&bentity.Id,
			&bentity.Type,
			&bentity.Title,
			&bentity.X,
			&bentity.Y,
			&bentity.FloorId,
			&bentity.Width,
			&bentity.Height,
			&bentity.Capacity,
			&bentity.CreatedAt,
			&bentity.UpdatedAt,
		)

		if err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		orders = append(orders, &entity.OrderBooking{Order: &order, Booking: &bentity})
	}

	return orders, nil
}

func (o *Order) GetById(c ctx.Context, id string) (*entity.Order, e.Error) {
	query, args, _ := sq.Select("*").From(orderTable).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := o.postgres.QueryRow(c, query, args...)

	var order entity.Order

	if err := order.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Order wan`t found.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &order, nil
}

func (o *Order) GetStats(c ctx.Context, filter string) (int, e.Error) {
	builder := sq.Select("*").From(orderTable)

	cur := time.Now()

	if filter == "day" {
		builder = builder.Where(sq.And{
			sq.GtOrEq{
				"created_at": cur.Add(time.Duration(-24) * time.Hour),
			},
			sq.LtOrEq{
				"created_at": cur,
			},
		})
	} else if filter == "week" {
		builder = builder.Where(sq.And{
			sq.GtOrEq{
				"created_at": cur.Add(time.Duration(-7*24) * time.Hour),
			},
			sq.LtOrEq{
				"created_at": cur,
			},
		})
	} else if filter == "month" {
		builder = builder.Where(sq.And{
			sq.GtOrEq{
				"created_at": cur.Add(time.Duration(-30*24) * time.Hour),
			},
			sq.LtOrEq{
				"created_at": cur,
			},
		})
	}

	query, args, _ := builder.PlaceholderFormat(sq.Dollar).ToSql()

	rows, err := o.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return 0, nil
		} else {
			return 0, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	count := 0

	for rows.Next() {
		var booking entity.Order

		if err := booking.Scan(rows); err != nil {
			return 0, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		count += 1
	}

	return count, nil
}
