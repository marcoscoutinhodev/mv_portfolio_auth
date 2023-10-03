package factory

import (
	"database/sql"
	"errors"
	"os"

	"github.com/marcoscoutinhodev/mv_chat/internal/domain/user"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/adapter"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/http/controller"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/repository"
)

type Auth struct {
	controller.Auth
}

func NewAuth(db *sql.DB) *Auth {
	hasher := adapter.NewHasher()
	repository := repository.NewUserRepository(db)
	queue := adapter.NewQueue()

	secretKey := os.Getenv("JWT_SECRET_KEY_DEFAULT")
	if secretKey == "" {
		panic(errors.New("JWT_SECRET_KEY_DEFAULT is not found in .env"))
	}
	encrypter := adapter.NewEncrypter(secretKey)

	usecase := user.NewUseCase(hasher, repository, queue, encrypter)
	controller := controller.NewAuth(*usecase)

	return &Auth{
		*controller,
	}
}
