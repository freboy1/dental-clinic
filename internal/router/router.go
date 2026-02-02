package router

import (
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"

	"dental_clinic/internal/config"
	"dental_clinic/internal/middleware"
	"dental_clinic/internal/modules/user"

	"github.com/gorilla/mux"
	gorilla_handler "github.com/gorilla/handlers"
	httpSwagger "github.com/swaggo/http-swagger"

)

func NewRouter(cfg *config.Config, db *pgxpool.Pool) http.Handler {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api").Subrouter()

	// Public routes
	public := api.NewRoute().Subrouter()
	user.RegisterPublicRoutes(public, db, cfg)

	// Private routes
	private := api.NewRoute().Subrouter()
	private.Use(middleware.JWTAuth(cfg.JWTSecret))
	user.RegisterPrivateRoutes(private, db, cfg)

	headersOk := gorilla_handler.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := gorilla_handler.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := gorilla_handler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	return gorilla_handler.CORS(originsOk, headersOk, methodsOk)(router)
}
