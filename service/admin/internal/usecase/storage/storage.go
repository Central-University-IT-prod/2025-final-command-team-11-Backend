package storage

import (
	"github.com/nikitaSstepanov/tools"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/sl"
	"REDACTED/team-11/backend/admin/internal/usecase/storage/booking"
	"REDACTED/team-11/backend/admin/internal/usecase/storage/booking_entity"
	"REDACTED/team-11/backend/admin/internal/usecase/storage/guest"
	"REDACTED/team-11/backend/admin/internal/usecase/storage/order"
	"REDACTED/team-11/backend/admin/internal/usecase/storage/verification"
	"REDACTED/team-11/backend/admin/pkg/client/minio"
)

type Storage struct {
	BookingEntity *booking_entity.BookingEntity
	Verification  *verification.Verification
	Booking       *booking.Booking
	Guest         *guest.Guest
	Order         *order.Order
	pg            pg.Client
	mn            minio.Client
}

type Config struct {
	Minio minio.Config `yaml:"minio"`
}

func New(c ctx.Context, cfg *Config) *Storage {
	pg := connectPg(c)
	minio := connectMinio(c, &cfg.Minio)

	return &Storage{
		BookingEntity: booking_entity.New(pg),
		Verification:  verification.New(pg, minio, cfg.Minio.Bucket),
		Booking:       booking.New(pg),
		Order:         order.New(pg),
		Guest:         guest.New(pg),
		pg:            pg,
		mn:            minio,
	}
}

func (s *Storage) Close() {
	s.pg.Close()
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

func connectMinio(c ctx.Context, cfg *minio.Config) minio.Client {
	log := c.Logger()

	minio, err := minio.New(c, cfg)
	if err != nil {
		log.Error("Can`t connect to minio.", sl.ErrAttr(err))
		panic("App start error.")
	} else {
		log.Info("Minio is connected.")
	}

	return minio
}
