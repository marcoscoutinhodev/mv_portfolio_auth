package factory

import (
	"errors"
	"os"

	"github.com/marcoscoutinhodev/mv_chat/internal/infra/adapter"
	"github.com/marcoscoutinhodev/mv_chat/internal/infra/http/middleware"
)

func NewMiddleware() *middleware.Middleware {
	secretKey := os.Getenv("JWT_SECRET_KEY_DEFAULT")
	if secretKey == "" {
		panic(errors.New("JWT_SECRET_KEY_DEFAULT is not found in .env"))
	}
	encrypter := adapter.NewEncrypter(secretKey)
	middleware := middleware.NewMiddleware(*encrypter)

	return middleware
}
