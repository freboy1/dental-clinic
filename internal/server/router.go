package server

import (
	"github.com/freboy1/dental-clinic/internal/middleware"
	"github.com/freboy1/dental-clinic/internal/utils"
	"net/http"

	"github.com/freboy1/dental-clinic/internal/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler, jwtManager *utils.JWTManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	protected := middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You are authorized ✅"))
	}))
	mux.Handle("/api/protected", protected)

	return mux
}
