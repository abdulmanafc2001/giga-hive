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

// Signup registers a new freelancer.
// @Summary Register a new freelancer
// @Description Register a new freelancer with the provided information.
// @Tags freelancer
// @Accept json
// @Produce json
// @Param freelancer body models.Freelancer true "freelancer registration information"
// @Success 200 {json} SuccessfulResponse "User registration successful"
// @Failure 400 {json} ErrorResponse "Bad request"
// @Failure 409 {json} ErrorResponse "Conflict - Username or phone number already exists"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /freelancer/signup [post]
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

	var freelancer models.Freelancer
	database.DB.Where("email = ?", input.Email).First(&freelancer)
	if freelancer.Email == input.Email {
		c.JSON(400, gin.H{
			"error": "This email already exist",
		})
		return
	}
	// checking the username already exist in database
	database.DB.Where("user_name = ?", input.User_Name).First(&freelancer)
	if freelancer.User_Name == input.User_Name {
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
		Phone:         input.Phone,
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

func ValidateOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get data",
		})
		return
	}

	var freelancer models.Freelancer
	if err := database.DB.Where("email = ? AND otp = ?", input.Email, input.OTP).First(&freelancer).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Incorrect user name and otp try again",
		})
		return
	}
	if err := database.DB.Model(&models.Freelancer{}).Where("email = ?", input.Email).Update("validate", true).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to update validate",
		})
		return
	}
	c.JSON(200, gin.H{"success": "Successfully validate freelancer"})
}
