package activation_code

import (
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type Code struct {
	redis rs.Client
}

func New(redis rs.Client) *Code {
	return &Code{
		redis,
	}
}

func (c *Code) Get(ctx ctx.Context, userId uint64) (*entity.ActivationCode, e.Error) {
	var result entity.ActivationCode

	err := c.redis.Get(ctx, redisKey(userId)).Scan(&result)
	if err != nil {
		if err == rs.Nil {
			return nil, notFoundErr.
				WithErr(err).
				WithCtx(ctx)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(ctx)
		}
	}

	return &result, nil
}

func (c *Code) Set(ctx ctx.Context, code *entity.ActivationCode) e.Error {
	err := c.redis.Set(ctx, redisKey(code.UserId), code, redisExpires).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}

func (c *Code) Del(ctx ctx.Context, userId uint64) e.Error {
	err := c.redis.Del(ctx, redisKey(userId)).Err()
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(ctx)
	}

	return nil
}
