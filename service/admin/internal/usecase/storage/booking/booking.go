package booking

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Booking struct {
	postgres pg.Client
}

func New(postgers pg.Client) *Booking {
	return &Booking{
		postgres: postgers,
	}
}

func (b *Booking) GetById(c ctx.Context, id string) (*entity.Booking, e.Error) {
	query, args, _ := sq.Select("*").From(bookingTable).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := b.postgres.QueryRow(c, query, args...)

	var booking entity.Booking

	if err := booking.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Book wadn`t found.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &booking, nil
}

func (b *Booking) GetNearest(c ctx.Context, id string) (*entity.Booking, e.Error) {
	cur := time.Now().UTC()

	builder := sq.Select("*").From(bookingTable).Where(sq.And{
		sq.Eq{
			"user_id": id,
		},
		sq.Or{
			sq.And{
				sq.GtOrEq{
					"time_from": cur,
				},
				sq.LtOrEq{
					"time_from": cur.Add(12 * time.Hour),
				},
			},
			sq.And{
				sq.LtOrEq{
					"time_from": cur,
				},
				sq.GtOrEq{
					"time_to": cur,
				},
			},
		},
	})

	query, args, _ := builder.
		OrderBy("time_from").
		PlaceholderFormat(sq.Dollar).ToSql()

	rows, err := b.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Where aren`t bookings in nearest 12 hours.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	bookings := make([]*entity.Booking, 0)

	for rows.Next() {
		var booking entity.Booking

		if err := booking.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		bookings = append(bookings, &booking)
	}

	if len(bookings) == 0 {
		return nil, e.New("Where aren`t bookings in nearest 12 hours.", e.NotFound).
			WithCtx(c)
	}

	return bookings[0], nil
}

func (b *Booking) GetNearestGuest(c ctx.Context, id string) (*entity.Booking, e.Error) {
	cur := time.Now().UTC()

	builder := sq.Select("*").From(fmt.Sprintf("%s AS b", bookingTable)).
		LeftJoin(fmt.Sprintf("%s AS g ON b.id = g.booking_id", guestTable)).
		Where(sq.And{
			sq.Eq{
				"g.user_id": id,
			},
			sq.Or{
				sq.And{
					sq.GtOrEq{
						"time_from": cur,
					},
					sq.LtOrEq{
						"time_from": cur.Add(12 * time.Hour),
					},
				},
				sq.And{
					sq.LtOrEq{
						"time_from": cur,
					},
					sq.GtOrEq{
						"time_to": cur,
					},
				},
			},
		})

	query, args, _ := builder.
		OrderBy("time_from").
		PlaceholderFormat(sq.Dollar).ToSql()

	rows, err := b.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Where aren`t invitations in nearest 12 hours.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	bookings := make([]*entity.Booking, 0)

	for rows.Next() {
		var booking entity.Booking
		var guest entity.Guest

		err := rows.Scan(
			&booking.Id,
			&booking.EntityId,
			&booking.UserId,
			&booking.TimeFrom,
			&booking.TimeTo,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&guest.UserId,
			&guest.BookingId,
			&guest.CreatedAt,
		)
		if err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		bookings = append(bookings, &booking)
	}

	if len(bookings) == 0 {
		return nil, e.New("Where aren`t invitations in nearest 12 hours.", e.NotFound).
			WithCtx(c)
	}

	return bookings[0], nil
}

func (b *Booking) GetStats(c ctx.Context, filter string) (int, e.Error) {
	builder := sq.Select("*").From(bookingTable)

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

	rows, err := b.postgres.Query(c, query, args...)
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
		var booking entity.Booking

		if err := booking.Scan(rows); err != nil {
			return 0, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		count += 1
	}

	return count, nil
}
