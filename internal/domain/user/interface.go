package user

import (
	"context"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Output struct {
	StatusCode int         `json:"-"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

type UseCaseInterface interface {
	Register(ctx context.Context, input *RegisterInput) *Output
}

type Repository interface {
	Find(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Store(ctx context.Context, user *entity.User, fn func() error) error
}

type Hasher interface {
	Generate(plaintext string) (string, error)
	Compare(hash, plaintext string) error
}

type Queue interface {
	RegisterNotification(ctx context.Context, user *entity.User) error
}
