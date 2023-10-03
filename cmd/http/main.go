package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	chi_middleware "github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/marcoscoutinhodev/mv_chat/factory"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_URI"))
	if err != nil {
		panic(err)
	}

	maxIdleConnsAsString := os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(maxIdleConnsAsString)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(maxIdleConns)

	auth := factory.NewAuth(db)
	middleware := factory.NewMiddleware()

	mux := chi.NewMux()
	mux.Use(chi_middleware.Logger)

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth.SignUp)
		r.Post("/signin", auth.SignIn)
		r.Post("/forgot-password", auth.ForgottenPassword)
		r.Post("/update-password", func(w http.ResponseWriter, r *http.Request) {
			middleware.Authorization(w, r, auth.UpdatePassword)
		})
	})

	if err := http.ListenAndServe(os.Getenv("SERVER_PORT"), mux); err != nil {
		log.Fatal(err)
	}
}
