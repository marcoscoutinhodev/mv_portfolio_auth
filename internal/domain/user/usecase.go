package user

import (
	"context"
	"net/http"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
)

type UseCase struct {
	hasher     Hasher
	repository Repository
	queue      Queue
	encrypter  Encrypter
}

func NewUseCase(hasher Hasher, repository Repository, queue Queue, encrypter Encrypter) *UseCase {
	return &UseCase{
		hasher:     hasher,
		repository: repository,
		queue:      queue,
		encrypter:  encrypter,
	}
}

func (s UseCase) Register(ctx context.Context, input *RegisterInput) (*Output, error) {
	u, err := s.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if u != nil {
		return &Output{
			StatusCode: http.StatusConflict,
			Error:      "email is already registered",
		}, nil
	}

	hashedPassword, err := s.hasher.Generate(input.Password)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &Output{
		StatusCode: http.StatusCreated,
		Data:       "check your inbox to verify your email and activate your account",
	}, nil
}

func (u UseCase) Auth(ctx context.Context, input *AuthInput) (*Output, error) {
	user, err := u.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if err := u.hasher.Compare(user.Password, input.Password); err == nil {
			token, refreshToken, err := u.encrypter.Encrypt(map[string]string{
				"sub": user.ID,
			}, 15, true)
			if err != nil {
				return nil, err
			}

			return &Output{
				StatusCode: http.StatusOK,
				Data: map[string]string{
					"accessToken":  token,
					"refreshToken": refreshToken,
				},
			}, nil
		}
	}

	return &Output{
		StatusCode: http.StatusUnauthorized,
		Error:      "invalid credentials",
	}, nil
}

func (u UseCase) ForgottenPassword(ctx context.Context, input *ForgottenPasswordInput) (*Output, error) {
	user, err := u.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		token, _, err := u.encrypter.Encrypt(map[string]string{
			"sub": user.ID,
		}, 60, false)
		if err != nil {
			return nil, err
		}

		if err := u.queue.ForgottenPasswordNotification(ctx, user, token); err != nil {
			return nil, err
		}
	}

	return &Output{
		StatusCode: http.StatusOK,
		Data:       "if the email provided is found, you will receive instructions to recover the password in your inbox",
	}, nil
}
