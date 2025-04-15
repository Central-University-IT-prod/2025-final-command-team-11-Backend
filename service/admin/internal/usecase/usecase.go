package usecase

import (
	"github.com/nikitaSstepanov/tools/httper"
	"REDACTED/team-11/backend/admin/internal/usecase/id"
	"REDACTED/team-11/backend/admin/internal/usecase/pkg/auth"
	"REDACTED/team-11/backend/admin/internal/usecase/pkg/booking"
	"REDACTED/team-11/backend/admin/internal/usecase/pkg/booking_entity"
	"REDACTED/team-11/backend/admin/internal/usecase/pkg/guest"
	"REDACTED/team-11/backend/admin/internal/usecase/pkg/order"
	"REDACTED/team-11/backend/admin/internal/usecase/pkg/verification"
	"REDACTED/team-11/backend/admin/internal/usecase/storage"
)

type UseCase struct {
	BookingEntity *booking_entity.BookingEntity
	Verification  *verification.Verification
	Booking       *booking.Booking
	Guest         *guest.Guest
	Order         *order.Order
	Auth          *auth.Auth
}

type Config struct {
	Jwt      auth.JwtOptions  `yaml:"jwt"`
	CoffeeId httper.ClientCfg `yaml:"coffee_id"`
}

func New(store *storage.Storage, cfg *Config) *UseCase {
	coffeeId := id.New(&cfg.CoffeeId)

	return &UseCase{
		Booking:       booking.New(store.Booking),
		BookingEntity: booking_entity.New(store.BookingEntity),
		Verification:  verification.New(store.Verification, coffeeId),
		Order:         order.New(store.Order),
		Auth:          auth.New(&cfg.Jwt),
		Guest:         guest.New(store.Guest, store.BookingEntity, store.Booking, coffeeId),
	}
}
