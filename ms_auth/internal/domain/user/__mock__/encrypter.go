package __mock__

import "github.com/stretchr/testify/mock"

type Encrypter struct {
	mock.Mock
}

func (m *Encrypter) Encrypt(payload map[string]interface{}, minutesToExpireToken uint) (token string, refreshToken string, err error) {
	args := m.Called(payload, minutesToExpireToken)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *Encrypter) Decrypt(token string) (payload map[string]interface{}, err error) {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *Encrypter) EncryptTemporary(payload map[string]interface{}) (string, error) {
	args := m.Called(payload)
	return args.String(0), args.Error(1)
}

func (m *Encrypter) DecryptTemporary(token string) (payload map[string]interface{}, err error) {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
