package auth

import (
	"auth_service/internal/config"
	"auth_service/internal/dto"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=Repository
type Repository interface {
	CreateUserSession(ctx context.Context, req dto.UserSessionRepoCreate) error
	GetUserSession(ctx context.Context, guid string, createdAt time.Time) (dto.UserSessionRepoGet, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=Hasher
type Hasher interface {
	GenerateFromPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=UUIDCreator
type UUIDCreator interface {
	New() string
}

type UseCase struct {
	repo        Repository
	hasher      Hasher
	uuidCreator UUIDCreator
	jwtSecret   string
}

func NewUseCase(repo Repository, hasher Hasher, uuidCreator UUIDCreator, jwtSecret string) UseCase {
	return UseCase{
		repo:        repo,
		hasher:      hasher,
		uuidCreator: uuidCreator,
		jwtSecret:   jwtSecret,
	}
}

func (u UseCase) CreateAccessToken(req dto.TokenCreateRequest) (string, error) {
	expTime := req.Time.Add(config.AccessTokenExpiration)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid":    req.GUID,
		"user_ip": req.UserIP,
		"iat":     req.Time.Unix(),
		"exp":     expTime.Unix(),
	})

	accessToken, err := claims.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return accessToken, nil
}

func (u UseCase) CreateRefreshToken(ctx context.Context, req dto.TokenCreateRequest) (string, error) {
	expTime := req.Time.Add(config.RefreshTokenExpiration)

	refreshToken := req.UserIP + u.uuidCreator.New()

	refreshBcrypt, err := u.hasher.GenerateFromPassword(refreshToken)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании хеша токена: %w", err)
	}

	repoReq := dto.UserSessionRepoCreate{
		GUID:         req.GUID,
		UserIP:       req.UserIP,
		RefreshToken: refreshBcrypt,
		CreatedAt:    req.Time,
		ExpiresAt:    expTime,
	}

	if err := u.repo.CreateUserSession(ctx, repoReq); err != nil {
		return "", fmt.Errorf("ошибка при обновлении токена в БД: %w", err)
	}

	return refreshToken, nil
}

func (u UseCase) CheckTokens(ctx context.Context, req dto.TokensCheckRequest) error {
	accessClaims, err := u.verifyToken(req.AccessToken)
	if err != nil {
		return err
	}

	accessCreatedAt := time.Unix(int64(accessClaims["iat"].(float64)), 0).UTC()

	refreshRepo, err := u.repo.GetUserSession(ctx, req.GUID, accessCreatedAt)
	if err != nil {
		return err
	}

	if refreshRepo.ExpiresAt.Before(req.Time) {
		return fmt.Errorf("refresh токен истек")
	}

	err = u.hasher.CompareHashAndPassword(refreshRepo.RefreshToken, req.RefreshToken)
	if err != nil {
		return fmt.Errorf("ошибка при проверке токена: %w", err)
	}

	if accessCreatedAt != refreshRepo.CreatedAt {
		return fmt.Errorf("время создания токенов не совпадает")
	}

	if accessClaims["user_ip"] != refreshRepo.UserIP {
		log.Println("ip поменялся, отправка предупреждения на почту")
	}

	return nil
}

func (u UseCase) verifyToken(tokenString string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(*jwt.Token) (any, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке токена: %w", err)
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("некорректный метод подписи: %v", token.Header["alg"])
	}

	if !token.Valid {
		return nil, fmt.Errorf("некорректный токен")
	}

	return claims, nil
}
