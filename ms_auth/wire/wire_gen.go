// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/marcoscoutinhodev/ms_auth/internal/domain/user"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/adapter"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/http/controller"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/http/mw"
)

// Injectors from wire.go:

func NewAuthController(db *sql.DB) *controller.Auth {
	adapterHasher := adapter.NewHasher()
	repository := user.NewRepository(db)
	adapterEmailNotification := adapter.NewEmailNotification()
	adapterEncrypter := adapter.NewEncrypter()
	adapterIDGenerator := adapter.NewIDGenerator()
	useCase := user.NewUseCase(adapterHasher, repository, adapterEmailNotification, adapterEncrypter, adapterIDGenerator)
	auth := controller.NewAuth(useCase)
	return auth
}

func NewAuthMiddleware() *mw.AuthMiddleware {
	adapterEncrypter := adapter.NewEncrypter()
	authMiddleware := mw.NewAuthMiddleware(adapterEncrypter)
	return authMiddleware
}

// wire.go:

var (
	hasher            = wire.NewSet(adapter.NewHasher, wire.Bind(new(adapter.HasherInterface), new(*adapter.Hasher)))
	userRepository    = wire.NewSet(user.NewRepository, wire.Bind(new(user.RepositoryInterface), new(*user.Repository)))
	emailNotification = wire.NewSet(adapter.NewEmailNotification, wire.Bind(new(adapter.EmailNotificationInterface), new(*adapter.EmailNotification)))
	encrypter         = wire.NewSet(adapter.NewEncrypter, wire.Bind(new(adapter.EncrypterInterface), new(*adapter.Encrypter)))
	userUseCase       = wire.NewSet(user.NewUseCase, wire.Bind(new(user.UseCaseInterface), new(*user.UseCase)))
	idGenerator       = wire.NewSet(adapter.NewIDGenerator, wire.Bind(new(adapter.IDGeneratorInterface), new(*adapter.IDGenerator)))
)
