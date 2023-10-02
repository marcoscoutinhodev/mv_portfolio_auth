package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
)

type UseCase struct {
	hasher     Hasher
	repository Repository
	queue      Queue
}

func NewUseCase(hasher Hasher, repository Repository, queue Queue) *UseCase {
	return &UseCase{
		hasher:     hasher,
		repository: repository,
		queue:      queue,
	}
}

func (s UseCase) Register(ctx context.Context, input *RegisterInput) *Output {
	u, err := s.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		return &Output{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	if u != nil {
		return &Output{
			StatusCode: http.StatusConflict,
			Error:      errors.New("email is already registered"),
		}
	}

	hashedPassword, err := s.hasher.Generate(input.Password)
	if err != nil {
		return &Output{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	user := &entity.User{Name: input.Name, Email: input.Email, Password: hashedPassword}

	// this function will ensure that the user who has just been stored in the database
	// will receive the verification email, otherwise it will return an error and the
	// user will not be stored in the database
	if err := s.repository.Store(ctx, user, func() error {
		if err := s.queue.RegisterNotification(ctx, user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return &Output{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return &Output{
		StatusCode: http.StatusCreated,
		Data:       "check your inbox to verify your email and activate your account",
	}
}
