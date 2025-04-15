package booking_entity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type BookingEntity struct {
	postgres pg.Client
}

func New(postgres pg.Client) *BookingEntity {
	return &BookingEntity{
		postgres: postgres,
	}
}

func (b *BookingEntity) GetEntity(c ctx.Context, id string) (*entity.BookingEntity, e.Error) {
	query, args, _ := sq.Select("*").From(bookingTable).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := b.postgres.QueryRow(c, query, args...)

	var booking entity.BookingEntity

	if err := booking.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Entity not found.", e.NotFound).
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

func (b *BookingEntity) GetFloors(c ctx.Context) ([]*entity.FloorEntity, e.Error) {
	query, args, _ := sq.Select("*").From(floorTable).ToSql()

	rows, err := b.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Floors not found.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	floors := make([]*entity.FloorEntity, 0)

	for rows.Next() {
		var floor entity.FloorEntity

		if err := floor.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		floors = append(floors, &floor)
	}

	return floors, nil
}

func (b *BookingEntity) GetFloor(c ctx.Context, id string) (*entity.FloorEntity, e.Error) {
	query, args, _ := sq.Select("*").From(floorTable).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := b.postgres.QueryRow(c, query, args...)

	var floor entity.FloorEntity

	if err := floor.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Floor not found.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &floor, nil
}

func (b *BookingEntity) CreateFloor(c ctx.Context, floor *entity.FloorEntity) e.Error {
	query, args, _ := sq.Insert(floorTable).
		Columns(
			"id", "name", "created_at", "updated_at",
		).
		Values(
			floor.Id, floor.Name, floor.CreatedAt, floor.UpdatedAt,
		).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := b.postgres.Begin(c)
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

func (b *BookingEntity) UpdateFloor(c ctx.Context, floor *entity.FloorEntity) e.Error {
	query, args, _ := sq.Update(floorTable).
		Set("name", floor.Name).Set("updated_at", floor.UpdatedAt).
		Where(sq.Eq{"id": floor.Id}).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := b.postgres.Begin(c)
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

func (b *BookingEntity) UpdateEntity(c ctx.Context, ent *entity.BookingEntity) e.Error {
	query, args, _ := sq.Update(bookingTable).
		Set("title", ent.Title).Set("x", ent.X).
		Set("y", ent.Y).Set("width", ent.Width).Set("height", ent.Height).
		Where(sq.Eq{"id": ent.Id}).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := b.postgres.Begin(c)
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

func (b *BookingEntity) GetEntities(c ctx.Context, id string) ([]*entity.BookingEntity, e.Error) {
	query, args, _ := sq.Select("*").From(bookingTable).
		Where(sq.Eq{"floor_id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	rows, err := b.postgres.Query(c, query, args...)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("Floor doesn`t have entities.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	entities := make([]*entity.BookingEntity, 0)

	for rows.Next() {
		var entity entity.BookingEntity

		if err := entity.Scan(rows); err != nil {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}

		entities = append(entities, &entity)
	}

	return entities, nil
}

func (b *BookingEntity) CreateEntity(c ctx.Context, entity *entity.BookingEntity) e.Error {
	query, args, _ := sq.Insert(bookingTable).
		Columns(
			"id", "type", "title", "x", "y", "floor_id",
			"width", "height", "capacity", "created_at", "updated_at",
		).
		Values(
			entity.Id, entity.Type, entity.Title, entity.X,
			entity.Y, entity.FloorId, entity.Width,
			entity.Height, entity.Capacity, entity.CreatedAt, entity.UpdatedAt,
		).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := b.postgres.Begin(c)
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

func (b *BookingEntity) DeleteEntity(c ctx.Context, id string) e.Error {
	query, args, _ := sq.Delete(bookingTable).Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := b.postgres.Begin(c)
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
		e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

func (b *BookingEntity) DeleteFloor(c ctx.Context, id string) e.Error {
	query, args, _ := sq.Delete(floorTable).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := b.postgres.Begin(c)
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
		e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}
