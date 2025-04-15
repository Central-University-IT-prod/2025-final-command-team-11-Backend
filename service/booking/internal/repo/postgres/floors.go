package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"REDACTED/team-11/backend/booking/internal/models"
)

var (
	floorsTable = "entity_floor"
)

type FloorsRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewFloorsRepo(db *sqlx.DB) *FloorsRepo {
	return &FloorsRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (fr *FloorsRepo) GetById(ctx context.Context, id uuid.UUID) (models.Floor, error) {
	op := "postgres.FloorsRepo.GetById"

	query, args, err := fr.sq.
		Select("*").
		From(floorsTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.Floor{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var res models.Floor
	if err := fr.db.GetContext(ctx, &res, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Floor{}, models.ErrFloorNotFound
		}

		return models.Floor{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return res, nil
}
