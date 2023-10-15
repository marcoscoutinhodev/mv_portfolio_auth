package user

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/marcoscoutinhodev/ms_auth/internal/domain/user/__mock__"
	"github.com/marcoscoutinhodev/ms_auth/internal/entity"
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
	return entity.NewUser("any_id", "any_name", "any_email", "any_password")
}

func (s *RegisterSuite) TestGivenAnErrorOnFindByEmail_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
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
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
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
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestGivenAnErrorOnEncryptTemporary_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{
		"sub": s.UserMock().ID,
	}).Return("", errors.New("any_error"))

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestGivenAnErrorOnStoreRepository_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	userMock := s.UserMock()
	userMock.Password = "hashed_password"
	repositoryMock.On("Store", context.Background(), userMock, mock.AnythingOfType("func() error")).Return(errors.New("any_error"))

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{
		"sub": userMock.ID,
	}).Return("any_token", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *RegisterSuite) TestSuccessScenario() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, nil)

	userMock := s.UserMock()
	userMock.Password = "hashed_password"
	repositoryMock.On("Store", context.Background(), userMock, mock.AnythingOfType("func() error")).Return(nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{
		"sub": userMock.ID,
	}).Return("any_token", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.Register(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusCreated, output.StatusCode)
	assert.Equal(s.T(), "check your inbox to verify your email and activate your account", output.Data)
	assert.Nil(s.T(), output.Error)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
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
	u := entity.NewUser("any_id", "any_name", "any_email", "hashed_password")
	u.ConfirmedEmail = true

	return u
}

func (s *AuthSuite) TestGivenAnErrorOnFindByEmail_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
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
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
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
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
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
	encrypter.On("Encrypt", map[string]interface{}{
		"sub": "any_id",
	}, uint(10)).Return("", "", errors.New("any_error"))

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
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
	encrypter.On("Encrypt", map[string]interface{}{
		"sub": "any_id",
	}, uint(10)).Return("any_access_token", "any_refresh_token", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.Auth(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
		Data: map[string]string{
			"name":         "any_name",
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
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
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
	encrypter.On("EncryptTemporary", map[string]interface{}{"sub": "any_id"}).Return("", errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.ForgottenPassword(context.Background(), s.InputMock())

	assert.Nil(s.T(), output)
	assert.Equal(s.T(), errors.New("any_error"), err)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *ForgottenPasswordSuite) TestGivenAnErrorInForgottenPassword_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{"sub": "any_id"}).Return("any_token", nil)

	emailNotificationMock := __mock__.EmailNotificationMock{}
	emailNotificationMock.On("ForgottenPassword", context.Background(), s.UserMock(), "any_token").Return(errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&emailNotificationMock,
		&encrypter,
		&__mock__.IDGeneratorMock{},
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
	encrypter.On("EncryptTemporary", map[string]interface{}{"sub": "any_id"}).Return("any_token", nil)

	emailNotificationMock := __mock__.EmailNotificationMock{}
	emailNotificationMock.On("ForgottenPassword", context.Background(), s.UserMock(), "any_token").Return(nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&emailNotificationMock,
		&encrypter,
		&__mock__.IDGeneratorMock{},
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

type UpdatePasswordSuite struct {
	suite.Suite
}

func (s *UpdatePasswordSuite) InputMock() *UpdatePasswordInput {
	return &UpdatePasswordInput{
		UserID:   "any_id",
		Password: "any_password",
	}
}

func (s *UpdatePasswordSuite) UserMock() *entity.User {
	return entity.NewUser("any_id", "any_name", "any_email", "hashed_password_from_db")
}

func (s *UpdatePasswordSuite) TestGivenAnErrorOnFind_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("Find", context.Background(), "any_id").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.UpdatePassword(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *UpdatePasswordSuite) TestGivenAnNilUser_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("Find", context.Background(), "any_id").Return(nil, nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.UpdatePassword(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("user not found in database"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *UpdatePasswordSuite) TestGivenAnErrorInHasherGeneration_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("Find", context.Background(), "any_id").Return(s.UserMock(), nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("", errors.New("any_error"))

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.UpdatePassword(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func (s *UpdatePasswordSuite) TestGivenAnErrorInUpdate_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("Find", context.Background(), "any_id").Return(s.UserMock(), nil)
	userMock := s.UserMock()
	userMock.Password = "hashed_password"
	repositoryMock.On("Update", context.Background(), userMock).Return(errors.New("any_error"))

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.UpdatePassword(context.Background(), s.InputMock())

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func (s *UpdatePasswordSuite) TestShouldReturnOKOnSuccess() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("Find", context.Background(), "any_id").Return(s.UserMock(), nil)
	userMock := s.UserMock()
	userMock.Password = "hashed_password"
	repositoryMock.On("Update", context.Background(), userMock).Return(nil)

	hasherMock := __mock__.HasherMock{}
	hasherMock.On("Generate", "any_password").Return("hashed_password", nil)

	sut := NewUseCase(
		&hasherMock,
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.UpdatePassword(context.Background(), s.InputMock())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
	}, *output)

	repositoryMock.AssertExpectations(s.T())
	hasherMock.AssertExpectations(s.T())
}

func TestUpdatePasswordSuite(t *testing.T) {
	suite.Run(t, new(UpdatePasswordSuite))
}

type EmailConfirmationRequest struct {
	suite.Suite
}

func (s *EmailConfirmationRequest) UserMock() *entity.User {
	return entity.NewUser("any_id", "any_name", "any_email", "hashed_password_from_db")
}

func (s *EmailConfirmationRequest) TestGivenAnErrorOnFindByEmail_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(nil, errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.EmailConfirmationRequest(context.Background(), "any_email")

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *EmailConfirmationRequest) TestGivenAnErrorOnEncryptTemporary_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{
		"sub": s.UserMock().ID,
	}).Return("", errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.EmailConfirmationRequest(context.Background(), "any_email")

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
}

func (s *EmailConfirmationRequest) TestGivenAnErrorOnEmailNotificationRegister_ShouldReturnError() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{
		"sub": s.UserMock().ID,
	}).Return("any_token", nil)

	emailNotificationMock := __mock__.EmailNotificationMock{}
	emailNotificationMock.On("Register", context.Background(), s.UserMock(), "any_token").Return(errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&emailNotificationMock,
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.EmailConfirmationRequest(context.Background(), "any_email")

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
	emailNotificationMock.AssertExpectations(s.T())
}

func (s *EmailConfirmationRequest) TestShouldReturnOKOnSuccess() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("FindByEmail", context.Background(), "any_email").Return(s.UserMock(), nil)

	encrypter := __mock__.Encrypter{}
	encrypter.On("EncryptTemporary", map[string]interface{}{
		"sub": s.UserMock().ID,
	}).Return("any_token", nil)

	emailNotificationMock := __mock__.EmailNotificationMock{}
	emailNotificationMock.On("Register", context.Background(), s.UserMock(), "any_token").Return(nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&emailNotificationMock,
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.EmailConfirmationRequest(context.Background(), "any_email")

	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
		Data:       "if the email provided is found, you will receive instructions to confirm your email in your inbox",
	}, *output)
	assert.Nil(s.T(), err)

	repositoryMock.AssertExpectations(s.T())
	encrypter.AssertExpectations(s.T())
	emailNotificationMock.AssertExpectations(s.T())
}

func TestEmailConfirmationRequest(t *testing.T) {
	suite.Run(t, new(EmailConfirmationRequest))
}

type ConfirmEmailSuite struct {
	suite.Suite
}

func (s *ConfirmEmailSuite) TestShouldReturnErrorOnFail() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("ConfirmEmail", context.Background(), "any_id").Return(errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.ConfirmEmail(context.Background(), "any_id")

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	repositoryMock.AssertExpectations(s.T())
}

func (s *ConfirmEmailSuite) TestShouldReturnOKOnSuccess() {
	repositoryMock := __mock__.RepositoryMock{}
	repositoryMock.On("ConfirmEmail", context.Background(), "any_id").Return(nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&repositoryMock,
		&__mock__.EmailNotificationMock{},
		&__mock__.Encrypter{},
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.ConfirmEmail(context.Background(), "any_id")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
	}, *output)

	repositoryMock.AssertExpectations(s.T())
}

func TestConfirmEmailSuite(t *testing.T) {
	suite.Run(t, new(ConfirmEmailSuite))
}

type NewAccessTokenSuite struct {
	suite.Suite
}

func (s *NewAccessTokenSuite) TestShouldReturnErrorOnFail() {
	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]interface{}{
		"sub": "any_id",
	}, uint(10)).Return("", "", errors.New("any_error"))

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&__mock__.RepositoryMock{},
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.NewAccessToken(context.Background(), "any_id")

	assert.Equal(s.T(), errors.New("any_error"), err)
	assert.Nil(s.T(), output)

	encrypter.AssertExpectations(s.T())
}

func (s *NewAccessTokenSuite) TestShouldReturnOKOnSuccess() {
	encrypter := __mock__.Encrypter{}
	encrypter.On("Encrypt", map[string]interface{}{
		"sub": "any_id",
	}, uint(10)).Return("any_access_token", "any_refresh_token", nil)

	sut := NewUseCase(
		&__mock__.HasherMock{},
		&__mock__.RepositoryMock{},
		&__mock__.EmailNotificationMock{},
		&encrypter,
		&__mock__.IDGeneratorMock{},
	)

	output, err := sut.NewAccessToken(context.Background(), "any_id")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Output{
		StatusCode: http.StatusOK,
		Data: map[string]string{
			"accessToken":  "any_access_token",
			"refreshToken": "any_refresh_token",
		},
	}, *output)

	encrypter.AssertExpectations(s.T())
}

func TestNewAccessTokenSuite(t *testing.T) {
	suite.Run(t, new(NewAccessTokenSuite))
}
