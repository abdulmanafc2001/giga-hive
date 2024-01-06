package usercontrollers

import (
	"strconv"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// declaring a global variable for validation
var (
	validate = validator.New()
)

// This user signup handler function ------------------------------------------------------->

// Signup registers a new user.
// @Summary Register a new user
// @Description Register a new user with the provided information.
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User registration information"
// @Success 200 {json} SuccessfulResponse "User registration successful"
// @Failure 400 {json} ErrorResponse "Bad request"
// @Failure 409 {json} ErrorResponse "Conflict - Username or phone number already exists"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /user/signup [post]
func UserSignup(c *gin.Context) {
	var input models.User
	if err := c.Bind(&input); err != nil {
		resp := helpers.Response{
			StatusCode: 422,
			Err:        "failed to parse request body. Please ensure it's valid JSON",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// validating the struct with given validate package
	if err := validate.Struct(input); err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        err.Error(),
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
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
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Password must contain at least one number and one special character.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the email already have in database
	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.Email == input.Email {
		resp := helpers.Response{
			StatusCode: 409,
			Err:        "An account with this email already exists",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the username already exist in database
	database.DB.Where("user_name = ?", input.User_Name).First(&user)
	if user.User_Name == input.User_Name {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "An account with this username already exists",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// hashing user given password
	ps, err := helpers.HashPassword(input.Password)
	if err != nil {
		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem processing your request. Please try again later.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// getting otp
	otp := helpers.GenerateOtp()

	if err = helpers.SendOtp(strconv.Itoa(otp), input.Email); err != nil {
		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem sending the OTP. Please try again later",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// creating the new user
	if err := database.DB.Create(&models.User{
		First_Name: input.First_Name,
		Last_Name:  input.Last_Name,
		User_Name:  input.User_Name,
		Email:      input.Email,
		Phone:      input.Phone,
		Password:   string(ps),
		Otp:        otp,
	}).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem creating your account. Please try again later.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// sending the otp to specified email address

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "otp send to " + input.Email,
	}
	helpers.ResponseResult(c, resp)
}

type OtpVerifiaction struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

// otp verification function ------------------------------------------------------------------>

// Signup Otp verification.
// @Summary Otp verfication of a new user
// @Description Otp verification with email id.
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body OtpVerifiaction true "User registration information"
// @Success 200 {json} SuccessfulResponse "User registration successful"
// @Failure 400 {json} ErrorResponse "Bad request"
// @Failure 409 {json} ErrorResponse "Conflict - Username or phone number already exists"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /user/signup/otpverification [post]
func OtpVerification(c *gin.Context) {
	// otp and geting from user
	var otp OtpVerifiaction
	if err := c.Bind(&otp); err != nil {
		resp := helpers.Response{
			StatusCode: 422,
			Err:        "failed to parse request body. Please ensure it's valid JSON",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", otp.Email).First(&user).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 404,
			Err:        "User not found",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the user already validate or not.
	if user.Validate {
		resp := helpers.Response{
			StatusCode: 409,
			Err:        "User is already verified. Please log in directly",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the otp correct or not.
	if otp.Otp != user.Otp {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Invalid OTP entered. Please check your OTP and try again.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// if the otp is correct the value in database validate column is updating to true
	database.DB.Model(&models.User{}).Where("id = ?", user.Id).Update("validate", true)

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully created user",
	}
	helpers.ResponseResult(c, resp)
}

// login function ------------------------------------------------------------------------------>

// Login handles user authentication by validating credentials and generating a JWT token upon successful login.
// @Summary User login
// @Description Authenticate user with provided credentials and generate JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body models.Login true "User login information"
// @Success 200 {json} SuccessfulResponse "Login successful"
// @Failure 401 {json} ErrorResponse "Unauthorized - Incorrect username or password"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /user/login [post]
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

	var user models.User
	if err := database.DB.Where("user_name = ? OR email = ?", input.UserName, input.UserName).First(&user).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 401,
			Err:        "Incorrect Username and password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking user status if user is blocked they cannot login
	if user.IsBlocked {
		resp := helpers.Response{
			StatusCode: 403,
			Err:        "Your account has been blocked.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking password if password is wrong it will return error
	if err := helpers.CheckPassword(user.Password, input.Password); err != nil {
		resp := helpers.Response{
			StatusCode: 401,
			Err:        "Incorrect username and password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// Genetated token using jwt
	token, err := helpers.GenerateJWT(user.Id)
	if err != nil {
		resp := helpers.Response{
			StatusCode: 500,
			Err:        "There was a problem processing your request.",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	// seting token into browser
	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("user_token", token, 3600*24, "", "", false, true)

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       token,
	}
	helpers.ResponseResult(c, resp)
}

func Logout(c *gin.Context) {
	c.SetCookie("user_token", "", -1, "", "", false, true)
	c.JSON(200, gin.H{
		"success": "Successfully logout",
	})
}

type bid struct {
	Description  string `json:"description"`
	About        string `json:"about"`
	MinPrice     int    `json:"minprice"`
	MaxPrice     int    `json:"maxprice"`
	ExpectedDays string `json:"expecteddays"`
	EndDay       int    `json:"endday"`
}
