package yandex

import (
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/client/pg"
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
)

type Yandex struct {
	postgres pg.Client
	redis    rs.Client
}

func New(postgres pg.Client, redis rs.Client) *Yandex {
	return &Yandex{
		postgres,
		redis,
	}
}

func (y *Yandex) GetById(ctx ctx.Context, id string) (*entity.Yandex, e.Error) {
	var yndx entity.Yandex

	query := idQuery(id)

	row := y.postgres.QueryRow(ctx, query)

	if err := yndx.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	err := y.redis.Set(ctx, redisKey(yndx.UserId), &yndx, redisExpires).Err()
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return &yndx, nil
}

func (y *Yandex) GetByUserId(ctx ctx.Context, userId uint64) (*entity.Yandex, e.Error) {
	var yndx entity.Yandex

	err := y.redis.Get(ctx, redisKey(userId)).Scan(&yndx)
	if err != nil && err != rs.Nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if yndx.YandexId != "" {
		return &yndx, nil
	}

	query := userIdQuery(userId)

	row := y.postgres.QueryRow(ctx, query)

	if err := yndx.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	err = y.redis.Set(ctx, redisKey(userId), &yndx, redisExpires).Err()
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return &yndx, nil
}

func (y *Yandex) Create(ctx ctx.Context, yndx *entity.Yandex) e.Error {
	log := ctx.Logger()

	query := createQuery(yndx)

	tx, err := y.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if _, err := tx.Exec(ctx, query); err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			log.Warn("transaction failed to rollback", sl.ErrAttr(err))
		}

		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	err = y.redis.Set(ctx, redisKey(yndx.UserId), yndx, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}
