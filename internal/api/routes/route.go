package routes

import (
	"net/http"

	"github.com/hasanraj3100/fridge-inventory/internal/api/handlers"
	"github.com/hasanraj3100/fridge-inventory/internal/middleware"
	"github.com/hasanraj3100/fridge-inventory/internal/utils"
)

type Router struct {
	mux               *http.ServeMux
	authHandler       *handlers.AuthHandler
	fridgeItemHandler *handlers.FridgeItemHandler
	itemUsageHandler  *handlers.ItemUsageHandler
	jwtManager        *utils.JWTManager
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	fridgeItemHandler *handlers.FridgeItemHandler,
	itemUsageHandler *handlers.ItemUsageHandler,
	jwtManager *utils.JWTManager,
) *Router {
	return &Router{
		mux:               http.NewServeMux(),
		authHandler:       authHandler,
		fridgeItemHandler: fridgeItemHandler,
		jwtManager:        jwtManager,
		itemUsageHandler:  itemUsageHandler,
	}
}

func (r *Router) Setup() http.Handler {
	r.mux.Handle("POST /api/v1/auth/register", http.HandlerFunc(r.authHandler.Register))
	r.mux.Handle("POST /api/v1/auth/login", http.HandlerFunc(r.authHandler.Login))
	r.mux.Handle("GET /health", http.HandlerFunc(r.healthCheck))
	r.mux.Handle("POST /api/v1/items", http.HandlerFunc(r.fridgeItemHandler.AddItem))
	r.mux.Handle("GET /api/v1/items", http.HandlerFunc(r.fridgeItemHandler.GetByUserID))
	r.mux.Handle("PATCH /api/v1/items/{id}", http.HandlerFunc(r.fridgeItemHandler.UpdateItem))
	r.mux.Handle("DELETE /api/v1/items/{id}", http.HandlerFunc(r.fridgeItemHandler.DeleteItem))
	r.mux.Handle("POST /api/v1/usage", http.HandlerFunc(r.itemUsageHandler.CreateItemUsage))
	r.mux.Handle("GET /api/v1/usage", http.HandlerFunc(r.itemUsageHandler.GetItemUsageByUserID))
	r.mux.Handle("PATCH /api/v1/usage/{id}", http.HandlerFunc(r.itemUsageHandler.UpdateItemUsage))

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
