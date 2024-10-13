package crypto

import "golang.org/x/crypto/bcrypt"

func Encrypt(input string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)

	return string(hashedPassword), err
}
