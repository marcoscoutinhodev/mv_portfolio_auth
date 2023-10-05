package __mock__

import "github.com/stretchr/testify/mock"

type HasherMock struct {
	mock.Mock
}

func (m *HasherMock) Generate(plaintext string) (string, error) {
	args := m.Called(plaintext)
	return args.String(0), args.Error(1)
}

func (m *HasherMock) Compare(hash, plaintext string) error {
	args := m.Called(hash, plaintext)
	return args.Error(0)
}
