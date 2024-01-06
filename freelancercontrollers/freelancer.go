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
		resp := helpers.Response{
			StatusCode: 422,
			Err:        "failed to parse request body. Please ensure it's valid JSON",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	if err := validate.Struct(input); err != nil {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        err.Error(),
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var freelancer models.Freelancer
	database.DB.Where("email = ?", input.Email).First(&freelancer)
	if freelancer.Email == input.Email {
		resp := helpers.Response{
			StatusCode: 409,
			Err:        "The email address you entered is already associated with an existing account",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the username already exist in database
	database.DB.Where("user_name = ?", input.User_Name).First(&freelancer)
	if freelancer.User_Name == input.User_Name {
		resp := helpers.Response{
			StatusCode: 409,
			Err:        "The username address you entered is already associated with an existing account",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
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

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Password must have at least one number, one special character",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	pas, err := helpers.HashPassword(input.Password)
	if err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "An error occurred while processing the password. Please try again later",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	strOTP := strconv.Itoa(helpers.GenerateOtp())

	if err := helpers.SendOtp(strOTP, input.Email); err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem sending the email with the OTP. Please check your email address and try again.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
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

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem creating your freelancer account. Please try again later.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 201,
		Err:        nil,
		Data:       "OTP sent to your email address",
	}
	helpers.ResponseResult(c, resp)

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
		resp := helpers.Response{
			StatusCode: 422,
			Err:        "failed to parse request body. Please ensure it's valid JSON",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var freelancer models.Freelancer
	if err := database.DB.Where("email = ? AND otp = ?", input.Email, input.Otp).First(&freelancer).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 401,
			Err:        "The email and OTP combination is incorrect. Please check your email and OTP carefully and try again.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	if err := database.DB.Model(&models.Freelancer{}).Where("email = ?", input.Email).Update("validate", true).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem updating your account validation. Please try again later or contact support.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 201,
		Err:        nil,
		Data:       "Successfully validate freelancer",
	}
	helpers.ResponseResult(c, resp)
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
		resp := helpers.Response{
			StatusCode: 422,
			Err:        "failed to parse request body. Please ensure it's valid JSON",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	var freelancer models.Freelancer
	if err := database.DB.Where("email = ? OR user_name = ?", input.UserName, input.UserName).
		First(&freelancer).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 401,
			Err:        "The email or username you entered is incorrect. Please check your credentials and try again.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if !freelancer.Validate && freelancer.IsBlocked {

		resp := helpers.Response{
			StatusCode: 401,
			Err:        "Your account is currently blocked. Please contact support for assistance.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := helpers.CheckPassword(freelancer.Password, input.Password); err != nil {

		resp := helpers.Response{
			StatusCode: 401,
			Err:        "Incorrect username and password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	tokenString, err := helpers.GenerateJWT(freelancer.Id)
	if err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem processing your request. Please try again later",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("freelancer_jwt", tokenString, 3600*24*30, "", "", false, true)

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       tokenString,
	}
	helpers.ResponseResult(c, resp)
}

func Logout(c *gin.Context) {
	c.SetCookie("freelancer_jwt", "", -1, "", "", false, true)
	c.JSON(200, gin.H{
		"success": "Successfully logout freelancer",
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
// @Success 200 {json} SuccessfulResponse "Successfully updated password"
// @Failure 400 {json} ErrorResponse "Bad Request"
// @Router /freelancer/profile [get]
func GetProfile(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")

	id := frlncr.(models.Freelancer).Id

	var freelancer1 freelancer
	if err := database.DB.Table("freelancers").Select("full_name,user_name,email,phone,qualification,tools").Where("id = ?", id).Scan(&freelancer1).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem retrieving the freelancer data. Please try again later ",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       freelancer1,
	}
	helpers.ResponseResult(c, resp)
}

func ChangePassword(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")
	id := frlncr.(models.Freelancer).Id

	var input models.CPassword
	if err := c.Bind(&input); err != nil {
		resp := helpers.Response{
			StatusCode: 422,
			Err:        "failed to parse request body. Please ensure it's valid JSON",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	
	var freelancer models.Freelancer
	if err := database.DB.First(&freelancer, id).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "The freelancer you're looking for was not found.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := helpers.CheckPassword(freelancer.Password, input.OldPassword); err != nil {

		resp := helpers.Response{
			StatusCode: 401,
			Err:        "The old password you entered is incorrect. Please try again.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if input.NewPassword != input.ConfirmPassword {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "The passwords you entered don't match. Please check and try again.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	pswd, err := helpers.HashPassword(input.NewPassword)
	if err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem processing your request. Please try again later",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := database.DB.Model(&models.Freelancer{}).Where("id = ?", id).Update("password", string(pswd)).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem updating your password. Please try again later",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully updated your password",
	}
	helpers.ResponseResult(c, resp)
}
