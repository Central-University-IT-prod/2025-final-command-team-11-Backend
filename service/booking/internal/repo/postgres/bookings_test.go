package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"REDACTED/team-11/backend/booking/pkg/postgres"
)

var (
	entityId    = uuid.MustParse("29dca021-adb9-4feb-bf71-f3adb21c41ab")
	timeFromStr = "2020-01-01T00:00:00Z"
	timeToStr   = "2020-01-10T00:00:00Z"
)

func TestBookingsRepo(t *testing.T) {
	ctx := context.Background()

	timeFrom, err := time.Parse(time.RFC3339, timeFromStr)
	require.NoError(t, err)
	timeTo, err := time.Parse(time.RFC3339, timeToStr)
	require.NoError(t, err)

	cfg := postgres.Config{
		Host:     "127.0.0.1",
		Port:     5432,
		DB:       "booking",
		User:     "booking_service",
		Password: "secret",
	}

	t.Log("port", cfg.Port)

	db, err := postgres.Connect(ctx, cfg)
	require.NoError(t, err)

	bookingsRepo := NewBookingsRepo(db)

	res, err := bookingsRepo.ListIntersectedForEntity(ctx, entityId, timeFrom, timeTo)
	require.NoError(t, err)

	ids := []string{}
	for _, r := range res {
		ids = append(ids, r.Id.String())
	}

	t.Log("ids", ids)
}
