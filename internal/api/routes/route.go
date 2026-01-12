package routes

import (
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/handlers"
	"github.com/hasanraj3100/fridge-inventory/internal/middleware"
	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

type Router struct {
	mux         *http.ServeMux
	authHandler *handlers.AuthHandler
	jwtManager  *utils.JWTManager
}

func NewRouter(authHandler *handlers.AuthHandler, jwtManager *utils.JWTManager) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		authHandler: authHandler,
		jwtManager:  jwtManager,
	}
}

func (r *Router) Setup() http.Handler {
	r.mux.HandleFunc("/api/v1/auth/register", r.authHandler.Register)
	r.mux.HandleFunc("/api/v1/auth/login", r.authHandler.Login)
	r.mux.HandleFunc("/health", r.healthCheck)

	handler := middleware.Chain(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.CORS(),
		middleware.AuthWithConfig(middleware.AuthConfig{
			JWTManager: r.jwtManager,
			Skipper: func(r *http.Request) bool {
				return r.URL.Path == "/api/v1/auth/register" ||
					r.URL.Path == "/api/v1/auth/login" ||
					r.URL.Path == "/health"
			},
		}),
	)(r.mux)

	return handler
}

func (r *Router) healthCheck(w http.ResponseWriter, rw *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
