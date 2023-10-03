package adapter

import (
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
