package adapter

import (
	"context"

	"github.com/marcoscoutinhodev/ms_auth/internal/entity"
)

type HasherInterface interface {
	Generate(plaintext string) (string, error)
	Compare(hash, plaintext string) error
}

type EmailNotificationInterface interface {
	Register(ctx context.Context, user *entity.User) error
	ForgottenPassword(ctx context.Context, user *entity.User, token string) error
}

type EncrypterInterface interface {
	Encrypt(payloadInterface interface{}, minutesToExpire uint, rt bool) (token string, refreshToken string, err error)
	Decrypt(token string) (payload map[string]interface{}, err error)
}
