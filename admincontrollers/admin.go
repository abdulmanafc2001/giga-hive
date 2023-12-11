package admincontrollers

import (
	"net/http"
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
		c.JSON(500, gin.H{
			"error": "Failed to geting data",
		})
		return
	}

	username, password := os.Getenv("ADMIN"), os.Getenv("ADMIN_PASSWORD")

	if username != input.UserName || password != input.Password {
		c.JSON(400, gin.H{
			"error": "incorrect username and password",
		})
		return
	}
	token, err := helpers.GenerateJWT(0)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "token generating error",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_admin", token, 3600*24, "", "", false, true)

	c.JSON(200, gin.H{
		"success": "successfully login admin",
		"token":   token,
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
		c.JSON(400, gin.H{
			"error": "Failed to find users list",
		})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
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
		c.JSON(400, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	if user.IsBlocked {
		c.JSON(400, gin.H{
			"error": "user already blocked",
		})
		return
	}
	// updating user is_blocked to true
	if err := database.DB.Model(&models.User{}).Where("id = ?", user_id).
		Update("is_blocked", "true").Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to Block user",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": "Successfully blocked user",
	})
}

func UnBlockUser(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	// fetching user data for checking user already unblocked?
	user, err := helpers.FindUserById(user_id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	if !user.IsBlocked {
		c.JSON(400, gin.H{
			"error": "user already unblocked",
		})
		return
	}
	// updating user is_blocked to true
	if err := database.DB.Model(&models.User{}).Where("id = ?", user_id).
		Update("is_blocked", "false").Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to Block user",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": "Successfully unblocked user",
	})
}


