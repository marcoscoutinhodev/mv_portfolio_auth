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

	if os.Getenv("REGISTER_NOTIFICATION_QUEUE") == "" {
		panic(errors.New("REGISTER_NOTIFICATION_QUEUE is not found in .env"))
	}

	if os.Getenv("REGISTER_NOTIFICATION_KEY") == "" {
		panic(errors.New("REGISTER_NOTIFICATION_KEY is not found in .env"))
	}

	if os.Getenv("REGISTER_NOTIFICATION_EXCHANGE") == "" {
		panic(errors.New("REGISTER_NOTIFICATION_EXCHANGE is not found in .env"))
	}

	if os.Getenv("FORGOT_PASSWORD_NOTIFICATION_QUEUE") == "" {
		panic(errors.New("FORGOT_PASSWORD_NOTIFICATION_QUEUE is not found in .env"))
	}

	if os.Getenv("FORGOT_PASSWORD_NOTIFICATION_KEY") == "" {
		panic(errors.New("FORGOT_PASSWORD_NOTIFICATION_KEY is not found in .env"))
	}

	if os.Getenv("FORGOT_PASSWORD_NOTIFICATION_EXCHANGE") == "" {
		panic(errors.New("FORGOT_PASSWORD_NOTIFICATION_EXCHANGE is not found in .env"))
	}

	emailNotification := adapter.NewEmailNotification()

	secretKey := os.Getenv("JWT_SECRET_KEY_DEFAULT")
	if secretKey == "" {
		panic(errors.New("JWT_SECRET_KEY_DEFAULT is not found in .env"))
	}
	encrypter := adapter.NewEncrypter(secretKey)

	usecase := user.NewUseCase(hasher, repository, emailNotification, encrypter)
	controller := controller.NewAuth(*usecase)

	return &Auth{
		*controller,
	}
}
