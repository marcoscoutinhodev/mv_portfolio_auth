package __mock__

import (
	"context"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/stretchr/testify/mock"
)

type QueueMock struct {
	mock.Mock
}

func (m *QueueMock) RegisterNotification(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
