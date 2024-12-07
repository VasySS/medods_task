package dto

import "time"

type TokenCreateRequest struct {
	Time   time.Time
	GUID   string
	UserIP string
}

type UserSessionRepoCreate struct {
	GUID         string
	UserIP       string
	RefreshToken string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type UserSessionRepoGet struct {
	ID           int       `db:"id"`
	UserID       string    `db:"user_id"`
	UserIP       string    `db:"user_ip"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}

type TokensRefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokensCheckRequest struct {
	GUID         string
	Time         time.Time
	AccessToken  string
	RefreshToken string
}

type CreatedTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
