package crypto

import "golang.org/x/crypto/bcrypt"

func Encrypt(input string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func IsSame(inputPlainPassword, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(inputPlainPassword))

	return err == nil
}
