package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/marcoscoutinhodev/ms_auth/config"
	"github.com/marcoscoutinhodev/ms_auth/wire"
)

func init() {
	config.Load()

	db, err := sql.Open(config.DB_DRIVER, config.DB_URI)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	filePath, err := filepath.Abs("migration")
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filePath),
		config.DB_NAME, driver,
	)
	if err != nil {
		panic(err)
	}

	m.Up()
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
			mw.AuthorizationTemporary(w, r, auth.UpdatePassword)
		})
		r.Post("/confirm-email", func(w http.ResponseWriter, r *http.Request) {
			mw.AuthorizationTemporary(w, r, auth.ConfirmEmail)
		})
	})

	if err := http.ListenAndServe(config.SERVER_PORT, mux); err != nil {
		log.Fatal(err)
	}
}
