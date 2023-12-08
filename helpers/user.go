package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func CheckPassword(dbPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(inputPassword))
}
