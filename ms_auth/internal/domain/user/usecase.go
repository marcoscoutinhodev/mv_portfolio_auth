package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/marcoscoutinhodev/ms_auth/internal/entity"
	"github.com/marcoscoutinhodev/ms_auth/internal/infra/adapter"
)

type UseCase struct {
	hasher            adapter.HasherInterface
	repository        RepositoryInterface
	emailNotification adapter.EmailNotificationInterface
	encrypter         adapter.EncrypterInterface
	idGenerator       adapter.IDGeneratorInterface
}

func NewUseCase(hasher adapter.HasherInterface, repository RepositoryInterface, emailNotification adapter.EmailNotificationInterface, encrypter adapter.EncrypterInterface, idGenerator adapter.IDGeneratorInterface) *UseCase {
	return &UseCase{
		hasher:            hasher,
		repository:        repository,
		emailNotification: emailNotification,
		encrypter:         encrypter,
		idGenerator:       idGenerator,
	}
}

func (u UseCase) Register(ctx context.Context, input *RegisterInput) (*Output, error) {
	us, err := u.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if us != nil {
		return &Output{
			StatusCode: http.StatusConflict,
			Error:      "email is already registered",
		}, nil
	}

	hashedPassword, err := u.hasher.Generate(input.Password)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(u.idGenerator.Generate(), input.Name, input.Email, hashedPassword)

	token, err := u.encrypter.EncryptTemporary(map[string]interface{}{
		"sub": user.ID,
	})
	if err != nil {
		return nil, err
	}

	// this function will ensure that the user who has just been stored in the database
	// will receive the verification email, otherwise it will return an error and the
	// user will not be stored in the database
	if err := u.repository.Store(ctx, user, func() error {
		if err := u.emailNotification.Register(ctx, user, token); err != nil {
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
			if !user.ConfirmedEmail {
				return &Output{
					StatusCode: http.StatusForbidden,
					Error:      "email must be verified",
				}, nil
			}

			token, refreshToken, err := u.encrypter.Encrypt(map[string]interface{}{
				"sub": user.ID,
			}, 10)
			if err != nil {
				return nil, err
			}

			return &Output{
				StatusCode: http.StatusOK,
				Data: map[string]string{
					"name":         user.Name,
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
		token, err := u.encrypter.EncryptTemporary(map[string]interface{}{
			"sub": user.ID,
		})
		if err != nil {
			return nil, err
		}

		if err := u.emailNotification.ForgottenPassword(ctx, user, token); err != nil {
			return nil, err
		}
	}

	return &Output{
		StatusCode: http.StatusOK,
		Data:       "if the email provided is found, you will receive instructions to recover the password in your inbox",
	}, nil
}

func (u UseCase) UpdatePassword(ctx context.Context, input *UpdatePasswordInput) (*Output, error) {
	user, err := u.repository.Find(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found in database")
	}

	hashedPassword, err := u.hasher.Generate(input.Password)
	if err != nil {
		return nil, err
	}

	user.UpdatePassword(hashedPassword)

	if err := u.repository.Update(ctx, user); err != nil {
		return nil, err
	}

	return &Output{
		StatusCode: http.StatusOK,
	}, nil
}

func (u UseCase) EmailConfirmationRequest(ctx context.Context, email string) (*Output, error) {
	user, err := u.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		token, err := u.encrypter.EncryptTemporary(map[string]interface{}{
			"sub": user.ID,
		})
		if err != nil {
			return nil, err
		}

		if err := u.emailNotification.Register(ctx, user, token); err != nil {
			return nil, err
		}
	}

	return &Output{
		StatusCode: http.StatusOK,
		Data:       "if the email provided is found, you will receive instructions to confirm your email in your inbox",
	}, nil
}

func (u UseCase) ConfirmEmail(ctx context.Context, userID string) (*Output, error) {
	if err := u.repository.ConfirmEmail(ctx, userID); err != nil {
		return nil, err
	}

	return &Output{
		StatusCode: http.StatusOK,
	}, nil
}

func (u UseCase) NewAccessToken(ctx context.Context, userID string) (*Output, error) {
	token, refreshToken, err := u.encrypter.Encrypt(map[string]interface{}{
		"sub": userID,
	}, 10)
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
