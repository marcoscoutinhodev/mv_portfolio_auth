//go:build wireinject
// +build wireinject

package wire

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/marcoscoutinhodev/ms_auth/internal/domain/user"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/adapter"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/http/controller"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/http/mw"
)

var (
	hasher            = wire.NewSet(adapter.NewHasher, wire.Bind(new(adapter.HasherInterface), new(*adapter.Hasher)))
	userRepository    = wire.NewSet(user.NewRepository, wire.Bind(new(user.RepositoryInterface), new(*user.Repository)))
	emailNotification = wire.NewSet(adapter.NewEmailNotification, wire.Bind(new(adapter.EmailNotificationInterface), new(*adapter.EmailNotification)))
	encrypter         = wire.NewSet(adapter.NewEncrypter, wire.Bind(new(adapter.EncrypterInterface), new(*adapter.Encrypter)))
	userUseCase       = wire.NewSet(user.NewUseCase, wire.Bind(new(user.UseCaseInterface), new(*user.UseCase)))
	idGenerator       = wire.NewSet(adapter.NewIDGenerator, wire.Bind(new(adapter.IDGeneratorInterface), new(*adapter.IDGenerator)))
)

func NewAuthController(db *sql.DB) *controller.Auth {
	wire.Build(
		hasher,
		userRepository,
		emailNotification,
		encrypter,
		userUseCase,
		idGenerator,
		controller.NewAuth,
	)
	return &controller.Auth{}
}

func NewAuthMiddleware() *mw.AuthMiddleware {
	wire.Build(
		encrypter,
		mw.NewAuthMiddleware,
	)
	return &mw.AuthMiddleware{}
}
