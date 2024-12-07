package config

import "time"

// Можно вынести в отдельный файл или .env (для JWTSecret нужно), но для этого отдельных требований не было,
// сделал всё проще

const (
	AccessTokenExpiration  = 24 * time.Hour
	RefreshTokenExpiration = 31 * 24 * time.Hour
	JWTSecret              = "123"
	BackendAddr            = "0.0.0.0:8080"
)
