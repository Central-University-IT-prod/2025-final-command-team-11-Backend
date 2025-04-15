package activation_code

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/assert/v2"
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
	"github.com/redis/go-redis/v9"
)

func TestGet(t *testing.T) {
	repo, err := setupRepo()
	if err != nil {
		t.Errorf("Can`t setup code repo: %v", err)
	}

	ctx := ctx.New(sl.Default())

	code := &entity.ActivationCode{
		Code:   gofakeit.DigitN(6),
		UserId: gofakeit.Uint64(),
	}

	err = repo.Set(ctx, code)
	if err != nil {
		t.Errorf("Can`t set code for tests: %v", err)
	}

	tests := []struct {
		TestName    string
		Code        string
		UserId      uint64
		IsError     bool
		ErrorStatus e.StatusType
	}{
		{
			TestName: "Success",
			Code:     code.Code,
			UserId:   code.UserId,
			IsError:  false,
		},
		{
			TestName:    "Code not found",
			UserId:      code.UserId + 1,
			IsError:     true,
			ErrorStatus: e.NotFound,
		},

		//TODO: Add more test cases
	}

	for _, tc := range tests {
		t.Run(tc.TestName, func(t *testing.T) {
			code, err := repo.Get(ctx, tc.UserId)
			if err != nil {
				if tc.IsError {
					assert.Equal(t, tc.ErrorStatus, err.GetCode())
				} else {
					t.Errorf("Test failing: %v", err)
				}
			} else {
				if tc.IsError {
					t.Error("Test failing: expected not nil error")
				} else {
					assert.Equal(t, tc.Code, code.Code)
				}
			}
		})
	}
}

func TestSet(t *testing.T) {
	repo, err := setupRepo()
	if err != nil {
		t.Errorf("Can`t setup code repo: %v", err)
	}

	ctx := ctx.New(sl.Default())

	tests := []struct {
		TestName    string
		Code        string
		UserId      uint64
		IsError     bool
		ErrorStatus e.StatusType
	}{
		{
			TestName: "Success",
			Code:     gofakeit.DigitN(6),
			UserId:   gofakeit.Uint64(),
			IsError:  false,
		},

		//TODO: Add more test cases
	}

	for _, tc := range tests {
		t.Run(tc.TestName, func(t *testing.T) {
			code := &entity.ActivationCode{
				Code:   tc.Code,
				UserId: tc.UserId,
			}

			err := repo.Set(ctx, code)
			if err != nil {
				if tc.IsError {
					assert.Equal(t, tc.ErrorStatus, err.GetCode())
				} else {
					t.Errorf("Test failing: %v", err)
				}
			} else {
				if tc.IsError {
					t.Error("Test failing: expected not nil error")
				} else {
					code, err := repo.Get(ctx, tc.UserId)
					if err != nil {
						t.Errorf("Test failing: %v", err)
					}

					assert.Equal(t, tc.Code, code.Code)
				}
			}
		})
	}
}

func TestDel(t *testing.T) {
	repo, err := setupRepo()
	if err != nil {
		t.Errorf("Can`t setup code repo: %v", err)
	}

	ctx := ctx.New(sl.Default())

	tests := []struct {
		TestName    string
		Code        string
		UserId      uint64
		IsError     bool
		ErrorStatus e.StatusType
	}{
		{
			TestName: "Success",
			Code:     gofakeit.DigitN(6),
			UserId:   gofakeit.Uint64(),
			IsError:  false,
		},

		//TODO: Add more test cases
	}

	for _, tc := range tests {
		t.Run(tc.TestName, func(t *testing.T) {
			code := &entity.ActivationCode{
				Code:   tc.Code,
				UserId: tc.UserId,
			}

			err := repo.Set(ctx, code)
			if err != nil {
				t.Errorf("Test failing: %v", err)
			}

			err = repo.Del(ctx, tc.UserId)
			if err != nil {
				if tc.IsError {
					assert.Equal(t, tc.ErrorStatus, err.GetCode())
				} else {
					t.Errorf("Test failing: %v", err)
				}
			} else {
				if tc.IsError {
					t.Error("Test failing: expected not nil error")
				} else {
					_, err := repo.Get(ctx, tc.UserId)
					if err == nil || err.GetCode() != e.NotFound {
						t.Errorf("Test failing: %v", err)
					}
				}
			}
		})
	}
}

func TestRedisKey(t *testing.T) {
	userId := gofakeit.Uint64()
	key := redisKey(userId)
	excpected := fmt.Sprintf("codes:%d", userId)

	assert.Equal(t, excpected, key)
}

func setupRepo() (*Code, e.Error) {
	server, err := miniredis.Run()
	if err != nil {
		return nil, e.E(err)
	}

	rs := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	return New(*rs), nil
}
