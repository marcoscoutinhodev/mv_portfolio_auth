package __mock__

import (
	"context"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/stretchr/testify/mock"
)

type EmailNotificationMock struct {
	mock.Mock
}

func (m *EmailNotificationMock) Register(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *EmailNotificationMock) ForgottenPassword(ctx context.Context, user *entity.User, token string) error {
	args := m.Called(ctx, user, token)
	return args.Error(0)
}
