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
