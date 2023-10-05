package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth.SignUp)
		r.Post("/signin", auth.SignIn)
		r.Post("/forgot-password", auth.ForgottenPassword)
		r.Post("/update-password", func(w http.ResponseWriter, r *http.Request) {
			mw.Authorization(w, r, auth.UpdatePassword)
		})
	})

	if err := http.ListenAndServe(config.SERVER_PORT, mux); err != nil {
		log.Fatal(err)
	}
}
