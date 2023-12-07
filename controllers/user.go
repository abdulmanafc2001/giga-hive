package controllers

import (
	"net/http"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
// declaring a global variable
var (
	validate = validator.New()
)

// This user signup handler function ------------------------------------------------------->
func UserSignup(c *gin.Context) {
	var input models.User
	if err := c.Bind(&input); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to get data",
		})
		return
	}
	// validating the struct with given validate package
	if err := validate.Struct(input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	// password validation
	hasNumber := false
	hasSpecialChar := false
	for _, char := range input.Password {
		switch {
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^':
			hasSpecialChar = true
		}
	}

	if !hasNumber || !hasSpecialChar {
		c.JSON(400, gin.H{
			"error": "password must have one special charecter and number",
		})
		return
	}
	// checking the email already have in database
	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.Email == input.Email {
		c.JSON(400, gin.H{
			"error": "This email already exist",
		})
		return
	}
	// hashing user given password
	ps, err := helpers.HashPassword(input.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Password hashing error",
		})
		return
	}
	// creating the new user
	if err := database.DB.Create(&models.User{
		First_Name: input.First_Name,
		Last_Name:  input.Last_Name,
		User_Name:  input.User_Name,
		Email:      input.Email,
		Password:   string(ps),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add data",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": "successfully created user",
	})
}
