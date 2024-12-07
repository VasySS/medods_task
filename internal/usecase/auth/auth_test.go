package auth_test

import (
	"auth_service/internal/config"
	"auth_service/internal/dto"
	"auth_service/internal/usecase/auth"
	"auth_service/internal/usecase/auth/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_CreateAccessToken(t *testing.T) {
	t.Parallel()

	type args struct {
		req dto.TokenCreateRequest
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успех",
			args: args{
				req: dto.TokenCreateRequest{
					GUID:   "guid",
					UserIP: "1.1.1.1",
					Time:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			want: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDQxNTM2MDAsImd1aWQiOiJn" +
				"dWlkIiwiaWF0IjoxNzA0MDY3MjAwLCJ1c2VyX2lwIjoiMS4xLjEuMSJ9.KirZrZod-6Kq31o_39k-gh" +
				"d03VhKaFigkw2TPS2qE4hwroB4O21YNLAIXp2C_cy3EA9Z-32-mcXZXeQgO-E11w",
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			hasher := mocks.NewHasher(t)
			u := auth.NewUseCase(repo, hasher, auth.GoogleUUID{}, "jwtsecret")

			got, err := u.CreateAccessToken(tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUseCase_CreateRefreshToken(t *testing.T) {
	t.Parallel()

	type args struct {
		req dto.TokenCreateRequest
	}

	tests := []struct {
		name    string
		args    args
		setup   func(m *mocks.Repository, hm *mocks.Hasher, um *mocks.UUIDCreator)
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успех",
			args: args{
				req: dto.TokenCreateRequest{
					GUID:   "user_guid",
					UserIP: "1.1.1.1",
					Time:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			setup: func(m *mocks.Repository, hm *mocks.Hasher, um *mocks.UUIDCreator) {
				hm.On("GenerateFromPassword", "MS4xLjEuMWd1aWQ=").
					Return("bcrypt_hash", nil)

				m.On("CreateUserSession", context.Background(), dto.UserSessionRepoCreate{
					GUID:         "user_guid",
					UserIP:       "1.1.1.1",
					RefreshToken: "bcrypt_hash",
					CreatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					ExpiresAt:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				}).Return(nil)

				um.On("New").Return("guid")
			},
			want:    "MS4xLjEuMWd1aWQ=",
			wantErr: assert.NoError,
		},
		{
			name: "ошибка хеширования",
			args: args{
				req: dto.TokenCreateRequest{
					GUID:   "user_guid",
					UserIP: "1.1.1.1",
					Time:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			setup: func(_ *mocks.Repository, hm *mocks.Hasher, um *mocks.UUIDCreator) {
				hm.On("GenerateFromPassword", "MS4xLjEuMWd1aWQ=").
					Return("", errors.New("ошибка хеширования"))

				um.On("New").Return("guid")
			},
			wantErr: assert.Error,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			hasher := mocks.NewHasher(t)
			uuidCreator := mocks.NewUUIDCreator(t)

			tt.setup(repo, hasher, uuidCreator)

			u := auth.NewUseCase(repo, hasher, uuidCreator, "jwtsecret")

			got, err := u.CreateRefreshToken(context.Background(), tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUseCase_CheckTokens(t *testing.T) {
	t.Parallel()

	testJWTSecret := "jwtsecret"
	creationTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iat": creationTime.Unix(),
	})
	accessToken, _ := accessClaims.SignedString([]byte(testJWTSecret))

	refreshReq := dto.TokensCheckRequest{
		GUID:         "guid",
		Time:         creationTime,
		AccessToken:  accessToken,
		RefreshToken: "refreshToken",
	}

	tests := []struct {
		name    string
		req     dto.TokensCheckRequest
		setup   func(*mocks.Repository, *mocks.Hasher)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успех",
			req:  refreshReq,
			setup: func(m *mocks.Repository, hm *mocks.Hasher) {
				m.On("GetUserSession", context.Background(), refreshReq.GUID, creationTime).
					Return(dto.UserSessionRepoGet{
						ID:           123,
						UserID:       refreshReq.GUID,
						UserIP:       "1.1.1.1",
						RefreshToken: "refreshTokenHash",
						CreatedAt:    creationTime,
						ExpiresAt:    creationTime.Add(config.RefreshTokenExpiration),
					}, nil)
				hm.On("CompareHashAndPassword", "refreshTokenHash", refreshReq.RefreshToken).
					Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "истек срок действия refresh token",
			req:  refreshReq,
			setup: func(m *mocks.Repository, _ *mocks.Hasher) {
				m.On("GetUserSession", context.Background(), refreshReq.GUID, refreshReq.Time).
					Return(dto.UserSessionRepoGet{
						ID:           123,
						UserID:       refreshReq.GUID,
						UserIP:       "1.1.1.1",
						RefreshToken: "refreshTokenHash",
						CreatedAt:    time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
						ExpiresAt:    time.Date(2023, 11, 2, 0, 0, 0, 0, time.UTC),
					}, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "ошибка хеширования",
			req:  refreshReq,
			setup: func(m *mocks.Repository, hm *mocks.Hasher) {
				m.On("GetUserSession", context.Background(), refreshReq.GUID, creationTime).
					Return(dto.UserSessionRepoGet{
						ExpiresAt: refreshReq.Time.Add(time.Hour),
					}, nil)

				hm.On("CompareHashAndPassword", "", refreshReq.RefreshToken).
					Return(errors.New("ошибка хеширования"))
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			hasher := mocks.NewHasher(t)
			uuidCreator := mocks.NewUUIDCreator(t)

			tt.setup(repo, hasher)

			u := auth.NewUseCase(repo, hasher, uuidCreator, testJWTSecret)

			tt.wantErr(t, u.CheckTokens(context.Background(), tt.req))
		})
	}
}
