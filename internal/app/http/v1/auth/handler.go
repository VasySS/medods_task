package auth

import (
	"auth_service/internal/dto"
	"context"
)

type UseCase interface {
	CreateAccessToken(req dto.TokenCreateRequest) (string, error)
	CreateRefreshToken(ctx context.Context, req dto.TokenCreateRequest) (string, error)
	CheckTokens(ctx context.Context, req dto.TokensCheckRequest) error
}

type Handler struct {
	usecase UseCase
}

func NewHandler(usecase UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
