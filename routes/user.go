package routes

import (
	"github.com/abdulmanafc2001/gigahive/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	r := router.Group("/user")
	r.POST("/signup",controllers.UserSignup)
}
