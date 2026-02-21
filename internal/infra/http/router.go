package infrahttp

import (
	"encoding/json"
	"health-checker/internal/infra/http/handlers"
	"health-checker/internal/infra/http/middlewares"
	"net/http"
)

type AppRouter struct {
	authHandler *handlers.AuthHandler
}

func NewAppRouter(authHandler *handlers.AuthHandler) *AppRouter {
	return &AppRouter{
		authHandler,
	}
}

func (r *AppRouter) InitializeRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", GetHealth)
	router.HandleFunc("POST /api/v1/auth/sign-up", r.authHandler.SignUp)
	router.HandleFunc("POST /api/v1/auth/login", r.authHandler.Login)
	router.Handle("POST /api/v1/auth/logout", middlewares.AuthMiddleware(http.HandlerFunc(r.authHandler.Logout)))
	return router
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	message := struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "App is healthy",
		Status:  "OK",
	}
	json.NewEncoder(w).Encode(message)
}
