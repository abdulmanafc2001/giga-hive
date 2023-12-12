package freelancercontrollers

import (
	"strconv"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func FreelancerSignup(c *gin.Context) {
	var input models.Freelancer
	if err := c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get data",
		})
		return
	}
	if err := validate.Struct(input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.Email == input.Email {
		c.JSON(400, gin.H{
			"error": "This email already exist",
		})
		return
	}
	// checking the username already exist in database
	database.DB.Where("user_name = ?", input.User_Name).First(&user)
	if user.User_Name == input.User_Name {
		c.JSON(400, gin.H{
			"error": "This username already exist",
		})
		return
	}

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

	strOTP := strconv.Itoa(helpers.GenerateOtp())

	if err := helpers.SendOtp(strOTP, input.Email); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to send email",
		})
		return
	}
	if err := database.DB.Create(&models.Freelancer{
		Full_Name:     input.Full_Name,
		User_Name:     input.User_Name,
		Email:         input.Email,
		Password:      input.Password,
		Qualification: input.Qualification,
		Tools:         input.Tools,
		OTP:           strOTP,
	}).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to create freelancer",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": "email sent to your email",
	})

}
