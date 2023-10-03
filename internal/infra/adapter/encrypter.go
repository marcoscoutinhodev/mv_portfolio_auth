package adapter

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Encrypter struct {
	secretKey string
}

func NewEncrypter(secretKey string) *Encrypter {
	return &Encrypter{
		secretKey: secretKey,
	}
}

func (e Encrypter) Encrypt(payload interface{}, minutesToExpire uint, rt bool) (token string, refreshToken string, err error) {
	claims := &jwt.MapClaims{
		"payload": payload,
		"registeredClaims": jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutesToExpire) * time.Minute)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(e.secretKey))
	if err != nil {
		return "", "", err
	}

	if rt {
		r := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": token,
			"registeredClaims": jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutesToExpire*6) * time.Minute)),
			},
		})
		refreshToken, err = r.SignedString([]byte(e.secretKey))
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
		return []byte(e.secretKey), nil
	}); err != nil {
		fmt.Println("to aq:", err)
		return nil, err
	}

	return claims, err
}
