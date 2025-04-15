package user

import (
	"fmt"

	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/client/pg"
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
)

type User struct {
	postgres pg.Client
	redis    rs.Client
}

func New(postgres pg.Client, redis rs.Client) *User {
	return &User{
		postgres,
		redis,
	}
}

func (u *User) Get(c ctx.Context, limit, offset int) ([]*entity.User, e.Error) {
	query := fmt.Sprintf("select * from users limit %d offset %d", limit, offset)

	rows, err := u.postgres.Query(c, query)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	users := make([]*entity.User, 0)

	for rows.Next() {
		var user entity.User

		if err := user.Scan(rows); err != nil {
			return nil, e.InternalErr.WithErr(err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *User) Getc(c ctx.Context) ([]*entity.User, e.Error) {
	query := "select * from users"

	rows, err := u.postgres.Query(c, query)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	users := make([]*entity.User, 0)

	for rows.Next() {
		var user entity.User

		if err := user.Scan(rows); err != nil {
			return nil, e.InternalErr.WithErr(err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *User) GetById(ctx ctx.Context, id string) (*entity.User, e.Error) {
	var user entity.User

	err := u.redis.Get(ctx, redisKey(id)).Scan(&user)
	if err != nil && err != rs.Nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if user.Id != "" {
		return &user, nil
	}

	query := idQuery(id)

	row := u.postgres.QueryRow(ctx, query)

	if err := user.Scan(row); err != nil {
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

	err = u.redis.Set(ctx, redisKey(id), &user, redisExpires).Err()
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return &user, nil
}

func (u *User) GetByEmail(ctx ctx.Context, email string) (*entity.User, e.Error) {
	var user entity.User

	query := emailQuery(email)

	row := u.postgres.QueryRow(ctx, query)

	if err := user.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		}

		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return &user, nil
}

func (u *User) Create(ctx ctx.Context, user *entity.User) e.Error {
	log := ctx.Logger()

	query := createQuery(user)

	tx, err := u.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	row := tx.QueryRow(ctx, query)

	err = row.Scan(&user.Id, &user.Role)
	if err != nil {
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

	err = u.redis.Set(ctx, redisKey(user.Id), user, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}

func (u *User) Update(ctx ctx.Context, user *entity.User) (*entity.User, e.Error) {
	log := ctx.Logger()

	query := updateQuery(user)

	tx, err := u.postgres.Begin(ctx)
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if _, err = tx.Exec(ctx, query); err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			log.Warn("transaction failed to rollback", sl.ErrAttr(err))
		}

		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	err = u.redis.Del(ctx, redisKey(user.Id)).Err()
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	user, getErr := u.GetById(ctx, user.Id)
	if getErr != nil {
		return nil, getErr
	}

	err = u.redis.Set(ctx, redisKey(user.Id), user, redisExpires).Err()
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return user, nil
}

func (u *User) AddRole(ctx ctx.Context, user *entity.User) e.Error {
	log := ctx.Logger()

	role := user.Role

	user, getErr := u.GetById(ctx, user.Id)
	if getErr != nil {
		return getErr
	}

	query := roleQuery(user, role)

	tx, err := u.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if _, err = tx.Exec(ctx, query); err != nil {
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

	err = u.redis.Del(ctx, redisKey(user.Id)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	user, getErr = u.GetById(ctx, user.Id)
	if getErr != nil {
		return getErr
	}

	err = u.redis.Set(ctx, redisKey(user.Id), user, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}

func (u *User) Verify(ctx ctx.Context, user *entity.User) e.Error {
	log := ctx.Logger()

	query := verifyQuery(user.Verified, user.Id)

	tx, err := u.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	if _, err = tx.Exec(ctx, query); err != nil {
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

	err = u.redis.Del(ctx, redisKey(user.Id)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	user, getErr := u.GetById(ctx, user.Id)
	if getErr != nil {
		return getErr.WithErr(err)
	}

	err = u.redis.Set(ctx, redisKey(user.Id), user, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}

func (u *User) Delete(ctx ctx.Context, user *entity.User) e.Error {
	log := ctx.Logger()

	query := deleteQuery()

	tx, err := u.postgres.Begin(ctx)
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	_, err = tx.Exec(ctx, query, user.Id)
	if err != nil {
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

	err = u.redis.Del(ctx, redisKey(user.Id)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}
