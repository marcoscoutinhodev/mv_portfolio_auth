package adapter

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/marcoscoutinhodev/mv_chat/config"
)

type Encrypter struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
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
	}
}

func (e Encrypter) generate(claims jwt.Claims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(e.privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (e Encrypter) Encrypt(payload interface{}, minutesToExpire uint, rt bool) (token string, refreshToken string, err error) {
	claims := &jwt.MapClaims{
		"payload": payload,
		"registeredClaims": jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutesToExpire) * time.Minute)),
		},
	}

	token, err = e.generate(claims)
	if err != nil {
		return "", "", err
	}

	if rt {
		claims := &jwt.MapClaims{
			"sub": token,
			"registeredClaims": jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutesToExpire*6) * time.Minute)),
			},
		}

		refreshToken, err = e.generate(claims)
		if err != nil {
			return "", "", err
		}

		return token, refreshToken, nil
	}

	return token, "", nil
}

func (e Encrypter) Decrypt(token string) (payload map[string]interface{}, err error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return e.publicKey, nil
	}); err != nil {
		fmt.Println("to aq:", err)
		return nil, err
	}

	return claims, err
}
