package verification

import (
	"bytes"

	sq "github.com/Masterminds/squirrel"
	"github.com/minio/minio-go/v7"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
	gominio "REDACTED/team-11/backend/admin/pkg/client/minio"
)

type Verification struct {
	postgres pg.Client
	minio    gominio.Client
	bucket   string
}

func New(posgtres pg.Client, minio gominio.Client, bucket string) *Verification {
	return &Verification{
		postgres: posgtres,
		minio:    minio,
		bucket:   bucket,
	}
}

func (v *Verification) CheckVerify(c ctx.Context, id string) (*entity.Verification, e.Error) {
	query, args, _ := sq.Select("*").From(verifyTable).
		Where(sq.Eq{"user_id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := v.postgres.QueryRow(c, query, args...)

	var data entity.Verification

	if err := data.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("User is not verified.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	data.PassportImage = "passport/" + data.PassportImage

	return &data, nil
}

func (v *Verification) Get(c ctx.Context, id string) (*entity.Verification, e.Error) {
	query, args, _ := sq.Select("*").From(verifyTable).
		Where(sq.Eq{"user_id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := v.postgres.QueryRow(c, query, args...)

	var data entity.Verification

	if err := data.Scan(row); err != nil {
		if err == pg.ErrNoRows {
			return nil, e.New("User is not verified.", e.NotFound).
				WithErr(err).
				WithCtx(c)
		} else {
			return nil, e.InternalErr.
				WithErr(err).
				WithCtx(c)
		}
	}

	return &data, nil
}

func (v *Verification) Verify(c ctx.Context, id string, image *entity.Image) e.Error {
	reader := bytes.NewReader(image.Buffer)

	_, err := v.minio.PutObject(c, v.bucket, image.Name, reader, image.Size, minio.PutObjectOptions{
		ContentType: image.ContentType,
	})

	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	query, args, _ := sq.Insert(verifyTable).
		Columns(
			"user_id", "passport_image",
		).
		Values(
			id, image.Name,
		).PlaceholderFormat(sq.Dollar).ToSql()

	tx, err := v.postgres.Begin(c)
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

func (v *Verification) UpdateData(c ctx.Context, image *entity.Image) e.Error {
	err := v.minio.RemoveObject(c, v.bucket, image.Name, minio.RemoveObjectOptions{})
	if err != nil {
		return e.InternalErr.WithErr(err).WithCtx(c)
	}

	reader := bytes.NewReader(image.Buffer)

	_, err = v.minio.PutObject(c, v.bucket, image.Name, reader, image.Size, minio.PutObjectOptions{
		ContentType: image.ContentType,
	})

	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}
