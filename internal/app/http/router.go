package http

import (
	"auth_service/internal/app/http/v1/auth"
	"auth_service/internal/config"
	"auth_service/internal/repository/postgres"
	authUsecase "auth_service/internal/usecase/auth"

	"github.com/go-chi/chi/v5"

	// swagger init
	_ "auth_service/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(repo *postgres.Repository) *chi.Mux {
	r := chi.NewRouter()

	au := authUsecase.NewUseCase(repo, authUsecase.BcryptHasher{}, authUsecase.GoogleUUID{}, config.JWTSecret)
	ah := auth.NewHandler(au)

	r.Route("/v1/auth/tokens", func(r chi.Router) {
		r.Get("/", ah.GenerateTokens)
		r.Post("/refresh", ah.RefreshTokens)
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(config.BackendAddr+"/swagger/doc.json")))

	return r
}
