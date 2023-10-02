package __mock__

import (
	"context"

	"github.com/marcoscoutinhodev/mv_chat/internal/entity"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Find(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *RepositoryMock) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *RepositoryMock) Store(ctx context.Context, user *entity.User, fn func() error) error {
	args := m.Called(ctx, user, fn)
	return args.Error(0)
}
