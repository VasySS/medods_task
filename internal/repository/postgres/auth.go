package postgres

import (
	"auth_service/internal/dto"
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r Repository) CreateUserSession(ctx context.Context, req dto.UserSessionRepoCreate) error {
	query := `
		INSERT INTO user_sessions
		(user_id, user_ip, refresh_token, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query, req.GUID, req.UserIP, req.RefreshToken, req.CreatedAt, req.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetUserSession(ctx context.Context, guid string, createdAt time.Time) (dto.UserSessionRepoGet, error) {
	query := `
		SELECT *
		FROM user_sessions
		WHERE user_id = $1 AND created_at = $2 
	`

	var resp dto.UserSessionRepoGet
	if err := pgxscan.Get(ctx, r.pool, &resp, query, guid, createdAt); err != nil {
		return resp, err
	}

	return resp, nil
}
