package admincontrollers

import (
	"net/http"
	"os"

	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

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
