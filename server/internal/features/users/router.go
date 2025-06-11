package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewRouter(db *pgxpool.Pool, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	userRepo := NewUserRepo(db)
	userService := NewUserService(userRepo)
	userHandler := NewHandler(userService, logger)

	r.Route("/users", func(r chi.Router) {
		r.Mount("/", userHandler.Routes())
	})

	return r
}
