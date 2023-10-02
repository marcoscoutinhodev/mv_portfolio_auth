package adapter

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
}

func (h Hasher) Generate(plaintext string) (string, error) {
	hasher, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hasher), nil
}

func (h Hasher) Compare(hash, plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
}
