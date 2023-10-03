package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth.SignUp)
		r.Post("/signin", auth.SignIn)
	})

	if err := http.ListenAndServe(os.Getenv("SERVER_PORT"), mux); err != nil {
		log.Fatal(err)
	}
}
