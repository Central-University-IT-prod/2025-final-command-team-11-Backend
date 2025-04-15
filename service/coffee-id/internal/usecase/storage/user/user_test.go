package user

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/assert/v2"
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	gopg "github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/redis/go-redis/v9"
)

var (
	YANDEX = string(types.YANDEX)
)

func TestGetById(t *testing.T) {
	rs, err := setupRedis()
	if err != nil {
		t.Errorf("Can`t setup user redis: %v", err)
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("Can`t setup pg mock: %v", err)
	}
	defer mock.Close()

	pg := gopg.NewWithMock(mock)

	repo := New(pg, *rs)

	ctx := ctx.New(sl.Default())

	tests := []struct {
		TestName    string
		Id          string
		Email       string
		Name        string
		Password    string
		Roles       types.Role
		Birthday    time.Time
		Verified    bool
		OAuth       *string
		IsError     bool
		ErrorMsg    string
		ErrorStatus e.StatusType
	}{
		{
			TestName: "Success",
			Id:       gofakeit.UUID(),
			Email:    gofakeit.Email(),
			Name:     gofakeit.Email(),
			Password: gofakeit.Letter(),
			Roles:    types.USER,
			Birthday: gofakeit.Date(),
			Verified: true,
			OAuth:    &YANDEX,
			IsError:  false,
		},
		{
			TestName:    "Not found",
			IsError:     true,
			Id:          gofakeit.UUID(),
			ErrorStatus: e.NotFound,
			ErrorMsg:    "This user wasn`t found.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.TestName, func(t *testing.T) {
			if !tc.IsError {
				rows := mock.NewRows([]string{"id", "email", "name", "password", "birthday", "roles", "verified", "oauth"}).
					AddRow(tc.Id, tc.Email, tc.Name, tc.Password, tc.Birthday, tc.Roles, tc.Verified, tc.OAuth)

				mock.ExpectQuery("SELECT * ").WillReturnRows(rows)
			} else {
				mock.ExpectQuery("SELECT *").WillReturnError(gopg.ErrNoRows)
			}

			user, err := repo.GetById(ctx, tc.Id)
			if err != nil {
				if tc.IsError {
					assert.Equal(t, tc.ErrorStatus, err.GetCode())
					assert.Equal(t, tc.ErrorMsg, err.GetMessage())
				} else {
					t.Errorf("Test failing: %v", err)
				}
			} else {
				if tc.IsError {
					t.Error("Test failing: expected not nil error")
				} else {
					assert.Equal(t, tc.Email, user.Email)
					assert.Equal(t, tc.Name, user.Name)
					assert.Equal(t, tc.Password, user.Password)
					assert.Equal(t, tc.Roles, user.Role)
					assert.Equal(t, tc.Birthday, user.Birthday)
					assert.Equal(t, tc.Verified, user.Verified)
					assert.Equal(t, *tc.OAuth, string(user.OAuth))
				}
			}
		})
	}
}

func TestIdQuery(t *testing.T) {
	id := gofakeit.UUID()
	query := idQuery(id)

	expected := fmt.Sprintf(
		`
			SELECT * FROM %s 
			WHERE id = '%s';
		`, usersTable, id,
	)

	assert.Equal(t, expected, query)
}

func TestEmailQuery(t *testing.T) {
	email := gofakeit.Email()
	query := emailQuery(email)

	expected := fmt.Sprintf(
		`
			SELECT * FROM %s 
			WHERE email = '%s';
		`, usersTable, email,
	)

	assert.Equal(t, expected, query)
}

func TestCreateQuery(t *testing.T) {
	user := &entity.User{
		Email:    gofakeit.Email(),
		Name:     gofakeit.Letter(),
		Password: gofakeit.Letter(),
		Birthday: gofakeit.Date(),
	}
	query := createQuery(user)

	expected := fmt.Sprintf(
		`
			INSERT INTO %s 
				(email, name, password, birthday) 
			VALUES 
				('%s', '%s', '%s', '%s') 
			RETURNING id, roles;
		`,
		usersTable, user.Email, user.Name, user.Password, user.Birthday.Format(time.DateOnly),
	)

	assert.Equal(t, expected, query)
}

func TestUpdateQuery(t *testing.T) {
	user := &entity.User{
		Id:       gofakeit.UUID(),
		Email:    gofakeit.Email(),
		Name:     gofakeit.FirstName(),
		Password: gofakeit.Password(true, true, true, true, true, 1),
		Birthday: gofakeit.Date(),
	}
	query := updateQuery(user)

	values := fmt.Sprintf(
		"email = '%s', name = '%s', password = '%s', birthday = '%s'",
		user.Email, user.Name, user.Password, user.Birthday.Format(time.DateOnly),
	)
	expected := fmt.Sprintf(
		`
			UPDATE %s 
			SET %s 
			WHERE id = '%s';
		`, usersTable, values, user.Id,
	)

	assert.Equal(t, expected, query)
}

func TestVerifyQuery(t *testing.T) {
	verify := gofakeit.Bool()
	id := gofakeit.UUID()
	query := verifyQuery(verify, id)

	expected := fmt.Sprintf(
		`
			UPDATE %s 
			SET verified = %t
			WHERE id = '%s';
		`, usersTable, verify, id,
	)

	assert.Equal(t, expected, query)
}

func TestDeleteQuery(t *testing.T) {
	query := deleteQuery()

	expected := fmt.Sprintf(
		`
			DELETE FROM %s 
			WHERE id = $1;
		`, usersTable,
	)

	assert.Equal(t, expected, query)
}

func TestRedisKey(t *testing.T) {
	userId := gofakeit.UUID()
	key := redisKey(userId)
	excpected := fmt.Sprintf("users:%s", userId)

	assert.Equal(t, excpected, key)
}

func TestSetupValues(t *testing.T) {
	user := &entity.User{
		Email:    gofakeit.Email(),
		Name:     gofakeit.FirstName(),
		Password: gofakeit.Password(true, true, true, true, true, 1),
		Birthday: gofakeit.Date(),
	}
	values := setupValues(user)

	expected := fmt.Sprintf(
		"email = '%s', name = '%s', password = '%s', birthday = '%s'",
		user.Email, user.Name, user.Password, user.Birthday.Format(time.DateOnly),
	)

	assert.Equal(t, values, expected)
}

func setupRedis() (*redis.Client, error) {
	server, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	rs := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	return rs, nil
}
