package auth

import (
	"auth_service/internal/dto"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// GenerateTokens генерирует новую пару access и refresh токенов.
// @Summary      Генерация новых JWT токенов
// @Description  Генерирует новую пару access и refresh токенов на основе GUID.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        guid  query     string  true  "Уникальный идентификатор (UUID)"
// @Success      200   {object}  dto.CreatedTokensResponse
// @Failure      400   {string}  string  "Некорректный GUID"
// @Failure      500   {string}  string  "Внутренняя ошибка сервера"
// @Router       /tokens [get]
func (h *Handler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if _, err := uuid.Parse(guid); err != nil {
		http.Error(w, "некорректный guid", http.StatusBadRequest)
		return
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	req := dto.TokenCreateRequest{
		Time:   time.Now().UTC(),
		GUID:   guid,
		UserIP: host,
	}

	h.createAndSendTokens(w, r, req)
}

// RefreshTokens обновляет access и refresh токены.
// @Summary      Обновление JWT токенов
// @Description  Обновляет access и refresh токены, если токены являются валидными.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        guid          query     string                 true  "Уникальный идентификатор (UUID)"
// @Param        request       body      dto.TokensRefreshRequest  true  "Access и Refresh токены"
// @Success      200           {object}  dto.CreatedTokensResponse
// @Failure      400           {string}  string  "Некорректный запрос или GUID"
// @Failure      401           {string}  string  "Некорректные токены"
// @Failure      500           {string}  string  "Внутренняя ошибка сервера"
// @Router       /tokens/refresh [post]
func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentTime := time.Now().UTC()

	guid := r.URL.Query().Get("guid")
	if _, err := uuid.Parse(guid); err != nil {
		http.Error(w, "некорректный guid", http.StatusBadRequest)
		return
	}

	var tokens dto.TokensRefreshRequest
	if err := render.DecodeJSON(r.Body, &tokens); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	checkReq := dto.TokensCheckRequest{
		GUID:         guid,
		Time:         currentTime,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	if err := h.usecase.CheckTokens(ctx, checkReq); err != nil {
		http.Error(w, "ошибка при проверке токенов", http.StatusUnauthorized)
		return
	}

	if err := h.usecase.SetSessionUsed(ctx, guid, currentTime); err != nil {
		http.Error(w, "ошибка при обновлении токена в БД", http.StatusInternalServerError)
		return
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	createReq := dto.TokenCreateRequest{
		Time:   currentTime,
		GUID:   guid,
		UserIP: host,
	}

	h.createAndSendTokens(w, r, createReq)
}

func (h *Handler) createAndSendTokens(w http.ResponseWriter, r *http.Request, createReq dto.TokenCreateRequest) {
	ctx := r.Context()

	accessToken, err := h.usecase.CreateAccessToken(createReq)
	if err != nil {
		http.Error(w, "ошибка при генерации токенов", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.usecase.CreateRefreshToken(ctx, createReq)
	if err != nil {
		http.Error(w, "ошибка при генерации токенов", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto.CreatedTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
