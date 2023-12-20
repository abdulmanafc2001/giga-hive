package admincontrollers

import (
	"os"
	"strconv"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

// Login handles admin authentication by validating credentials and generating a JWT token upon successful login.
// @Summary Admin login
// @Description Authenticate user with provided credentials and generate JWT token
// @Tags adminauthentication
// @Accept json
// @Produce json
// @Param user body models.Login true "admin login information"
// @Success 200 {json} SuccessfulResponse "Login successful"
// @Failure 401 {json} ErrorResponse "Unauthorized - Incorrect username or password"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /admin/login [post]
func Login(c *gin.Context) {
	var input models.Login
	if err := c.Bind(&input); err != nil {
		resp := helpers.Response{
			StatusCode: 500,
			Err:        "Failed to geting data",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	username, password := os.Getenv("ADMIN"), os.Getenv("ADMIN_PASSWORD")

	if username != input.UserName || password != input.Password {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Incorrect username and password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	token, err := helpers.GenerateJWT(0)
	if err != nil {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Token generating error",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       token,
	}
	helpers.ResponseResult(c, resp)
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt_admin", "", -1, "", "", false, true)
	c.JSON(200, gin.H{
		"success": "successfully logout admin",
	})
}

// ListUsers lists all users.
// @Summary List users
// @Description Retrieve a list of all users
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {json} UsersResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/user/list [get]
func GetUserList(c *gin.Context) {
	var users []models.UserList
	// fetching all users data from database
	if err := database.DB.Table("users").Select("id,first_name,last_name,user_name,email,phone,is_blocked,validate").Scan(&users).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find users list",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       users,
	}
	helpers.ResponseResult(c, resp)
}

// BlockUser blocks a user by updating the 'is_blocked' field in the database.
// @Summary Block a user
// @Description Block a user by updating the 'is_blocked' field in the database
// @Tags users
// @Accept json
// @Produce json
// @Param user_id query int true "User ID to block"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse "Successfully blocked user"
// @Failure 400 {json} ErrorResponse "Failed to find user" or "User already blocked" or "Failed to block user"
// @Router /admin/user/block [patch]
func BlockUser(c *gin.Context) {
	// geting userid
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	// fetching user data for checking user already blocked?
	user, err := helpers.FindUserById(user_id)
	if err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find user",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if user.IsBlocked {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "user already blocked",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// updating user is_blocked to true
	if err := database.DB.Model(&models.User{}).Where("id = ?", user_id).
		Update("is_blocked", "true").Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to Block user",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// c.JSON(200, gin.H{
	// 	"success": "Successfully blocked user",
	// })
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully blocked user",
	}
	helpers.ResponseResult(c, resp)

}

// BlockUser unblock a user by updating the 'is_blocked' field in the database.
// @Summary Unblock a user
// @Description Unblock a user by updating the 'is_blocked' field in the database
// @Tags users
// @Accept json
// @Produce json
// @Param user_id query int true "User ID to nblock"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse "Successfully unblocked user"
// @Failure 400 {json} ErrorResponse "Failed to find user" or "User already unblocked" or "Failed to unblock user"
// @Router /admin/user/unblock [patch]
func UnBlockUser(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	// fetching user data for checking user already unblocked?
	user, err := helpers.FindUserById(user_id)
	if err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find user",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)

		return
	}

	if !user.IsBlocked {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "user already unblocked",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// updating user is_blocked to true
	if err := database.DB.Model(&models.User{}).Where("id = ?", user_id).
		Update("is_blocked", "false").Error; err != nil {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to unblock user",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	// c.JSON(200, gin.H{
	// 	"success": "Successfully unblocked user",
	// })
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully unblocked user",
	}
	helpers.ResponseResult(c, resp)
}
