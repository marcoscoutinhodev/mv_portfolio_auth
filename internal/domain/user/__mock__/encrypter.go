package __mock__

import "github.com/stretchr/testify/mock"

type Encrypter struct {
	mock.Mock
}

func (m *Encrypter) Encrypt(payload interface{}, minutesToExpire uint, rt bool) (token string, refreshToken string, err error) {
	args := m.Called(payload, minutesToExpire, rt)
	return args.String(0), args.String(1), args.Error(2)
}
