package router

import (
	"encoding/json"
	"net/http"
)

type AppRouter struct{}

func NewAppRouter() *AppRouter {
	return &AppRouter{}
}

func (r *AppRouter) InitializeRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", GetHealth)
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
