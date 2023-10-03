package middleware

type Encrypter interface {
	Decrypt(token string) (payload map[string]interface{}, err error)
}
