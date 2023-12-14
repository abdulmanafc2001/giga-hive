package freelancercontrollers

import (
	"net/http"
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

	pas, err := helpers.HashPassword(input.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Password hashing error",
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
		Password:      string(pas),
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

type OtpVerifiaction struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

// Signup Otp verification.
// @Summary Otp verfication of a new freelancer
// @Description Otp verification with email id.
// @Tags freelancer
// @Accept json
// @Produce json
// @Param user body OtpVerifiaction true "User registration information"
// @Success 200 {json} SuccessfulResponse "User registration successful"
// @Failure 400 {json} ErrorResponse "Bad request"
// @Failure 409 {json} ErrorResponse "Conflict - Username or phone number already exists"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /freelancer/signup/otpverification [post]
func ValidateOTP(c *gin.Context) {
	var input OtpVerifiaction
	if err := c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get data",
		})
		return
	}

	var freelancer models.Freelancer
	if err := database.DB.Where("email = ? AND otp = ?", input.Email, input.Otp).First(&freelancer).Error; err != nil {
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

// Login handles freelancer authentication by validating credentials and generating a JWT token upon successful login.
// @Summary Freelancer login
// @Description Authenticate user with provided credentials and generate JWT token
// @Tags freelancer
// @Accept json
// @Produce json
// @Param user body models.Login true "User login information"
// @Success 200 {json} SuccessfulResponse "Login successful"
// @Failure 401 {json} ErrorResponse "Unauthorized - Incorrect username or password"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /freelancer/login [post]
func Login(c *gin.Context) {
	var input models.Login

	if err := c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get body",
		})
		return
	}
	var freelancer models.Freelancer
	if err := database.DB.Where("email = ? OR user_name = ?", input.UserName, input.UserName).
		First(&freelancer).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Incorrect username and password",
		})
		return
	}

	if !freelancer.Validate && freelancer.IsBlocked {
		c.JSON(401, gin.H{
			"error": "Unautharized access",
		})
		return
	}

	if err := helpers.CheckPassword(freelancer.Password, input.Password); err != nil {
		c.JSON(400, gin.H{
			"error": "Incorrect username and password",
		})
		return
	}

	tokenString, err := helpers.GenerateJWT(freelancer.Id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to generate JWT token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("freelancer_jwt", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{
		"message": "Successfully logged freelancer",
	})
}

type freelancer struct {
	Full_Name     string `json:"fullname" validate:"min=4,max=20"`
	User_Name     string `json:"username" validate:"min=4,max=20"`
	Email         string `json:"email" validate:"email"`
	Phone         string `json:"phone" validate:"min=10,max=10"`
	Qualification string `json:"qualification"`
	Tools         string `json:"tools"`
}

// @Summary Get freelancer profile
// @Description Get the profile of the logged-in freelancer
// @Tags Freelancer
// @Produce json
// @Param Authorization header string true "JWT Token" default(bearer <token>)
// @Success 200 {object} gin.H{"profile": models.Freelancer} "OK"
// @Failure 400 {object} gin.H{"error": string} "Bad Request"
// @Router /freelancer/profile [get]
func GetProfile(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")

	id := frlncr.(models.Freelancer).Id

	var freelancer1 freelancer
	if err := database.DB.Table("freelancers").Select("full_name,user_name,email,phone,qualification,tools").Where("id = ?", id).Scan(&freelancer1).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get data",
		})
		return
	}

	c.JSON(200, gin.H{
		"profile": freelancer1,
	})
}

func ChangePassword(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")
	id := frlncr.(models.Freelancer).Id

	var input models.CPassword
	if err := c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "failed to get data",
		})
		return
	}
	var freelancer models.Freelancer
	if err := database.DB.First(&freelancer, id).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	if err := helpers.CheckPassword(freelancer.Password, input.OldPassword); err != nil {
		c.JSON(400, gin.H{
			"error": "Incorrect old password",
		})
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		c.JSON(400, gin.H{
			"error": "Incorrect confirm password",
		})
		return
	}

	pswd, err := helpers.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if err := database.DB.Model(&models.Freelancer{}).Where("id = ?", id).Update("password", string(pswd)).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to update password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully updated password",
	})
}
