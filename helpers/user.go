package helpers

import (
	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func CheckPassword(dbPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(inputPassword))
}

func FindUserById(id int) (models.User, error) {
	var user models.User
	if err := database.DB.Find(&user, id).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
