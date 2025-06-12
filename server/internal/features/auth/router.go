package auth

import (
	"github.com/devmuggs/onigiri/server/internal/features/users"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewRouter(db *pgxpool.Pool, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	usersRepo := users.NewUserRepo(db)
	userService := users.NewUserService(usersRepo)
	handler := NewHandler(userService, logger)

	r.Route("/", func(r chi.Router) {
		r.Mount("/", handler.Routes())
	})

	return r
}
