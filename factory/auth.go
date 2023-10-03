package factory

import (
	"database/sql"

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
	usecase := user.NewUseCase(hasher, repository, queue)
	controller := controller.NewAuth(*usecase)

	return &Auth{
		*controller,
	}
}
