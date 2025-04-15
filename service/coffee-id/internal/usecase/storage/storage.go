package storage

import (
	code "github.com/nikitaSstepanov/coffee-id/internal/usecase/storage/activation_code"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/storage/user"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/storage/yandex"
	"github.com/nikitaSstepanov/tools"
	"github.com/nikitaSstepanov/tools/client/pg"
	rs "github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
)

type Storage struct {
	Users  *user.User
	Codes  *code.Code
	Yandex *yandex.Yandex
	pg     pg.Client
	rs     rs.Client
}

func New(c ctx.Context) *Storage {
	postgres := connectPg(c)
	redis := connectRs(c)

	return &Storage{
		Users:  user.New(postgres, redis),
		Codes:  code.New(redis),
		Yandex: yandex.New(postgres, redis),
		pg:     postgres,
		rs:     redis,
	}
}

func (s *Storage) Close() e.Error {
	s.pg.Close()
	return e.E(s.rs.Close())
}

func connectPg(c ctx.Context) pg.Client {
	log := c.Logger()

	postgres, err := tools.Pg()
	if err != nil {
		log.Error("Can`t connect to postgres.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Postgres is connected.")
	}

	if err := postgres.RegisterTypes(pgTypes); err != nil {
		log.Error("Can`t register custom types.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Custom types are registered.")
	}

	return postgres
}

func connectRs(c ctx.Context) rs.Client {
	log := c.Logger()

	redis, err := tools.Redis()
	if err != nil {
		log.Error("Can`t connect to redis.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Redis is connected.")
	}

	return redis
}
