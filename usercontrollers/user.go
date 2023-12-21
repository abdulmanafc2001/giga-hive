package usercontrollers

import (
	"strconv"
	"time"

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
		c.JSON(500, gin.H{
			"error": "Failed to get data",
		})
		return
	}
	// validating the struct with given validate package
	if err := validate.Struct(input); err != nil {
		// c.JSON(400, gin.H{
		// 	"error": err.Error(),
		// })
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
		// c.JSON(400, gin.H{
		// 	"error": "password must have one special charecter and number",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Password must have one special charecter and number",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the email already have in database
	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.Email == input.Email {
		// c.JSON(400, gin.H{
		// 	"error": "This email already exist",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "This email already exist",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the username already exist in database
	database.DB.Where("user_name = ?", input.User_Name).First(&user)
	if user.User_Name == input.User_Name {
		// c.JSON(400, gin.H{
		// 	"error": "This username already exist",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "This username already exist",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// hashing user given password
	ps, err := helpers.HashPassword(input.Password)
	if err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Password hashing error",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Password hashing error",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// getting otp
	otp := helpers.GenerateOtp()

	if err = helpers.SendOtp(strconv.Itoa(otp), input.Email); err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Failed to send otp",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to send otp",
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
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"error": "Failed to add data",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to add data",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// sending the otp to specified email address

	// c.JSON(200, gin.H{
	// 	"success": "otp send to " + input.Email,
	// })
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
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"error": "Failed to get data",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", otp.Email).First(&user).Error; err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Failed to find user",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find user",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the user already validate or not.
	if user.Validate {
		// c.JSON(400, gin.H{
		// 	"error": "User already verified",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "User already verified",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking the otp correct or not.
	if otp.Otp != user.Otp {
		// c.JSON(400, gin.H{
		// 	"error": "Otp verification failed, check your otp",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Otp verification failed, check your otp",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// if the otp is correct the value in database validate column is updating to true
	database.DB.Model(&models.User{}).Where("id = ?", user.Id).Update("validate", true)
	// c.JSON(200, gin.H{
	// 	"success": "successfully created user",
	// })
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
		// c.JSON(500, gin.H{
		// 	"error": "Find to get values",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var user models.User
	if err := database.DB.Where("user_name = ? OR email = ?", input.UserName, input.UserName).First(&user).Error; err != nil {
		// c.JSON(401, gin.H{
		// 	"error": "incorrect username and password",
		// })
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
		// c.JSON(400, gin.H{
		// 	"error": "User is blocked",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "User is blocked",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// checking password if password is wrong it will return error
	if err := helpers.CheckPassword(user.Password, input.Password); err != nil {
		// c.JSON(401, gin.H{
		// 	"error": "incorrect username and password",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Incorrect username and password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// Genetated token using jwt
	token, err := helpers.GenerateJWT(user.Id)
	if err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Failed to generate token",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to generate token",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	// seting token into browser
	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("user_token", token, 3600*24, "", "", false, true)

	// c.JSON(200, gin.H{
	// 	"success": "Login successfull",
	// 	"token":   token,
	// })
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

func CreateBid(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var input bid
	if err := c.Bind(&input); err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Failed to get body",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	endingDate := time.Now().Add(time.Hour * 24 * time.Duration(input.EndDay))
	endDate := endingDate.Format("2006-01-02")

	if err := database.DB.Create(&models.Bid{
		Description:  input.Description,
		About:        input.About,
		MinPrice:     input.MinPrice,
		MaxPrice:     input.MaxPrice,
		ExpectedDays: input.ExpectedDays,
		User_Id:      id,
		EndDay:       endDate,
	}).Error; err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Failed to create bid",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to create bid",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	// c.JSON(200, gin.H{
	// 	"success": "Successfully created new bid",
	// })
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully created new bid",
	}
	helpers.ResponseResult(c, resp)
}

type AuctionDetail struct {
	Id            int    `json:"id" gorm:"primaryKey"`
	BidId         int    `json:"bidid"`
	AuctionAmount string `json:"auctionamount"`
	Full_Name     string `json:"fullname"`
	Description   string `json:"description"`
	About         string `json:"about"`
	PriceRange    string `json:"pricerange"`
	ExpectedDays  string `json:"expecteddays"`
	EndDay        string `json:"endday"`
}

func GetAuctionedBid(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var auctions []AuctionDetail
	if err := database.DB.Table("auctions").Select("auctions.id,auctions.bid_id,auctions.auction_amount,freelancers.full_name,bids.description,bids.about,bids.price_range,bids.expected_days,bids.end_day").
		Joins("INNER JOIN freelancers ON freelancers.id=auctions.freelancer_id").
		Joins("INNER JOIN bids ON bids.id=auctions.bid_id").Where("bids.user_id=?", id).Scan(&auctions).Error; err != nil {
		// c.JSON(400, gin.H{
		// 	"error": "Failed to find auctioned datas",
		// })
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find auctioned datas",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	// c.JSON(200, gin.H{
	// 	"auctions": auctions,
	// })
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       auctions,
	}
	helpers.ResponseResult(c, resp)
}
