package user

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/marcoscoutinhodev/mv_chat/internal/domain/user/__mock__"
	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterSuite struct {
	suite.Suite
}

func (s *RegisterSuite) InputMock() *RegisterInput {
	return &RegisterInput{
		Name:     "any_name",
		Email:    "any_email",
		Password: "any_password",
	}
}

func (s *RegisterSuite) UserMock() *entity.User {
	return entity.NewUser("", "any_name", "any_email", "any_password")
}

func (s *RegisterSuite) TestGivenAnErrorOnFindByEmail_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestGivenAnEmailRegistered_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusConflict, output.StatusCode)
	assert.Equal(s.T(), "email is already registered", output.Error)
	assert.Nil(s.T(), output.Data)

	repositoryMock.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestGivenAnErrorInHasherGeneration_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("", errors.New("any_error"))

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestGivenAnErrorOnStoreRepository_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	userMock := s.UserMock()
	userMock.Password = "hashed_password"
	repositoryMock.On("Store", context.Background(), userMock, mock.AnythingOfType("func() error")).Return(errors.New("any_error"))

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestSuccessScenario() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	userMock := s.UserMock()
	userMock.Password = "hashed_password"
	repositoryMock.On("Store", context.Background(), userMock, mock.AnythingOfType("func() error")).Return(nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusCreated, output.StatusCode)
	assert.Equal(s.T(), "check your inbox to verify your email and activate your account", output.Data)
	assert.Nil(s.T(), output.Error)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

type AuthSuite struct {
	suite.Suite
}

func (s *AuthSuite) InputMock() *AuthInput {
	return &AuthInput{
		Email:    "any_email",
		Password: "any_password",
	}
}

func (s *AuthSuite) UserMock() *entity.User {
	return entity.NewUser("any_id", "any_name", "any_email", "hashed_password")
}

func (s *AuthSuite) TestGivenAnErrorOnFindByEmail_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Auth(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *AuthSuite) TestGivenAnErrorEmailNotRegister_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Auth(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusUnauthorized,
		Error:      "invalid credentials",
	}, *output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *AuthSuite) TestGivenAnInvalidPassword_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Compare", "hashed_password", "any_password").Return(errors.New("any_error"))

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.Auth(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusUnauthorized,
		Error:      "invalid credentials",
	}, *output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func (s *AuthSuite) TestGivenAnErrorInEncrypter_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Compare", "hashed_password", "any_password").Return(nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]string{
		"sub": "any_id",
	}, uint(15), true).Return("", "", errors.New("any_error"))

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.QueueMock{},
		&encrypter,
	)

	output, err := sut.Auth(context.Background(), s.InputMock())

	assert.Nil(s.T(), output)
	assert.Equal(s.T(), errors.New("any_error"), err)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *AuthSuite) TestGivenValidInput_ShouldReturnAccessTokenAndRefreshToken() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Compare", "hashed_password", "any_password").Return(nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]string{
		"sub": "any_id",
	}, uint(15), true).Return("any_access_token", "any_refresh_token", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.QueueMock{},
		&encrypter,
	)

	output, err := sut.Auth(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
		Data: map[string]string{
			"accessToken":  "any_access_token",
			"refreshToken": "any_refresh_token",
		},
	}, *output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

type ForgottenPasswordSuite struct {
	suite.Suite
}

func (s *ForgottenPasswordSuite) InputMock() *ForgottenPasswordInput {
	return &ForgottenPasswordInput{
		Email: "any_email",
	}
}

func (s *ForgottenPasswordSuite) UserMock() *entity.User {
	return entity.NewUser("any_id", "any_name", "any_email", "hashed_password")
}

func (s *ForgottenPasswordSuite) TestGivenAnErrorOnFindByEmail_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.QueueMock{},
		&__mock__.Encrypter{},
	)

	output, err := sut.ForgottenPassword(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *ForgottenPasswordSuite) TestGivenAnErrorInEncrypter_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]string{
		"sub": "any_id",
	}, uint(60), false).Return("", "", errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.QueueMock{},
		&encrypter,
	)

	output, err := sut.ForgottenPassword(context.Background(), s.InputMock())

	assert.Nil(s.T(), output)
	assert.Equal(s.T(), errors.New("any_error"), err)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *ForgottenPasswordSuite) TestGivenAnErrorInForgottenPasswordNotification_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]string{
		"sub": "any_id",
	}, uint(60), false).Return("any_token", "", nil)

	queueMock := __mock__.QueueMock{}
	queueMock.On("ForgottenPasswordNotification", context.Background(), s.UserMock(), "any_token").Return(errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&queueMock,
		&encrypter,
	)

	output, err := sut.ForgottenPassword(context.Background(), s.InputMock())

	assert.Nil(s.T(), output)
	assert.Equal(s.T(), errors.New("any_error"), err)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *ForgottenPasswordSuite) TestGivenNoError_ShouldReturnGenericResponse() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]string{
		"sub": "any_id",
	}, uint(60), false).Return("any_token", "", nil)

	queueMock := __mock__.QueueMock{}
	queueMock.On("ForgottenPasswordNotification", context.Background(), s.UserMock(), "any_token").Return(nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&queueMock,
		&encrypter,
	)

	output, err := sut.ForgottenPassword(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
		Data:       "if the email provided is found, you will receive instructions to recover the password in your inbox",
	}, *output)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func TestForgottenPasswordSuite(t *testing.T) {
	suite.Run(t, new(ForgottenPasswordSuite))
}
