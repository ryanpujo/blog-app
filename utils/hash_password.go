package utils

import "golang.org/x/crypto/bcrypt"

var HashPassword = encryptPassword

func encryptPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
