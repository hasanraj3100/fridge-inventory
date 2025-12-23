package routes

import (
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/handlers"
	"github.com/hasanraj3100/fridge-inventory/internal/middleware"
)

type Router struct {
	mux         *http.ServeMux
	authHandler *handlers.AuthHandler
}

func NewRouter(authHandler *handlers.AuthHandler) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		authHandler: authHandler,
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
	)(r.mux)

	return handler
}

func (r *Router) healthCheck(w http.ResponseWriter, rw *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
