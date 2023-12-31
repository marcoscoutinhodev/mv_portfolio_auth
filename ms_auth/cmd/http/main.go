package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/marcoscoutinhodev/ms_auth/config"
	"github.com/marcoscoutinhodev/ms_auth/wire"
)

func init() {
	config.Load()
}

func main() {
	db, err := sql.Open(config.DB_DRIVER, config.DB_URI)
	if err != nil {
		panic(err)
	}

	maxIdleConns, err := strconv.Atoi(config.DB_MAX_IDLE_CONNS)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(maxIdleConns)

	auth := wire.NewAuthController(db)
	mw := wire.NewAuthMiddleware()

	mux := chi.NewMux()
	mux.Use(middleware.Logger)
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedHeaders: []string{"*"},
	}))

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth.SignUp)
		r.Post("/signin", auth.SignIn)
		r.Post("/forgot-password", auth.ForgottenPassword)
		r.Post("/update-password", func(w http.ResponseWriter, r *http.Request) {
			mw.AuthorizationTemporary(w, r, auth.UpdatePassword)
		})
		r.Post("/email-confirmation-request", auth.EmailConfirmationRequest)
		r.Post("/confirm-email", func(w http.ResponseWriter, r *http.Request) {
			mw.AuthorizationTemporary(w, r, auth.ConfirmEmail)
		})
		r.Post("/refresh-token", func(w http.ResponseWriter, r *http.Request) {
			mw.Authorization(w, r, auth.NewAccessToken)
		})
	})

	if err := http.ListenAndServe(config.SERVER_PORT, mux); err != nil {
		log.Fatal(err)
	}
}
