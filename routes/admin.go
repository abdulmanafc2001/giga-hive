package routes

import (
	controllers "github.com/abdulmanafc2001/gigahive/admincontrollers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	r := router.Group("/admin")
	r.POST("/login",controllers.Login)
}