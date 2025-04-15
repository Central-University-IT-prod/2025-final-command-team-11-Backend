package guest

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Guest struct {
	postgres pg.Client
}

func New(postgres pg.Client) *Guest {
	return &Guest{
		postgres: postgres,
	}
}

func (g *Guest) Create(c ctx.Context, guest *entity.Guest) e.Error {
	query, args, _ := sq.Insert(guestTable).
		Columns(
			"user_id", "booking_id", "created_at",
		).
		Values(
			guest.UserId, guest.BookingId, guest.CreatedAt,
		).
		PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := g.postgres.Begin(c)
	if err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}
	defer tx.Rollback(c)

	if _, err := tx.Exec(c, query, args...); err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}

	if err := tx.Commit(c); err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}

	return nil
}

func (g *Guest) Get(c ctx.Context, id string) ([]*entity.Guest, e.Error) {
	query, args, _ := sq.Select("*").From(guestTable).
		Where(sq.Eq{"booking_id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	rows, err := g.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Booking hasn`t inites.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	guests := make([]*entity.Guest, 0)

	for rows.Next() {
		var guest entity.Guest

		if err := guest.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		guests = append(guests, &guest)
	}

	return guests, nil
}

func (g *Guest) GetById(c ctx.Context, bookId, userId string) (*entity.Guest, e.Error) {
	query, args, _ := sq.Select("*").From(guestTable).
		Where(sq.And{
			sq.Eq{
				"user_id": userId,
			},
			sq.Eq{
				"booking_id": bookId,
			},
		}).PlaceholderFormat(sq.Dollar).ToSql()

	row := g.postgres.QueryRow(c, query, args...)

	var guest entity.Guest

	if err := guest.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Invite was`nt found.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &guest, nil
}

func (g *Guest) Delete(c ctx.Context, bookId string, userId string) e.Error {
	query, args, _ := sq.Delete(guestTable).
		Where(sq.And{
			sq.Eq{
				"user_id": userId,
			},
			sq.Eq{
				"booking_id": bookId,
			},
		}).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := g.postgres.Begin(c)
	if err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}
	defer tx.Rollback(c)

	if _, err := tx.Exec(c, query, args...); err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}

	if err := tx.Commit(c); err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}

	return nil
}
