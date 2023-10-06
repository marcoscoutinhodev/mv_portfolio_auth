package adapter

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/marcoscoutinhodev/ms_auth/config"
)

type Encrypter struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	secretKey  string
}

func NewEncrypter() *Encrypter {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.JWT_PRIVATE_KEY))
	if err != nil {
		panic(err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.JWT_PUBLIC_KEY))
	if err != nil {
		panic(err)
	}
	return &Encrypter{
		privateKey: privateKey,
		publicKey:  publicKey,
		secretKey:  config.JWT_SECRET_KEY,
	}
}

func (e Encrypter) generate(claims jwt.Claims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(e.privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (e Encrypter) mapClaims(payload map[string]interface{}, minutesToExpireToken uint) *jwt.MapClaims {
	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Duration(minutesToExpireToken) * time.Minute))

	return &claims
}

func (e Encrypter) Encrypt(payload map[string]interface{}, minutesToExpireToken uint) (token string, refreshToken string, err error) {
	claims := e.mapClaims(payload, minutesToExpireToken)
	token, err = e.generate(claims)
	if err != nil {
		return "", "", err
	}

	claims = e.mapClaims(payload, minutesToExpireToken*12)
	refreshToken, err = e.generate(claims)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (e Encrypter) Decrypt(token string) (payload map[string]interface{}, err error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return e.publicKey, nil
	}); err != nil {
		return nil, err
	}

	return claims, err
}

func (e Encrypter) EncryptTemporary(payload map[string]interface{}) (string, error) {
	claims := e.mapClaims(payload, 60)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(e.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (e Encrypter) DecryptTemporary(token string) (payload map[string]interface{}, err error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(e.secretKey), nil
	}); err != nil {
		return nil, err
	}

	return claims, err
}
