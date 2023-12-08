package routes

import (
	controllers "github.com/abdulmanafc2001/gigahive/usercontrollers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	r := router.Group("/user")
	r.POST("/signup", controllers.UserSignup)
	r.POST("/signup/otpverification", controllers.OtpVerification)
	r.POST("/login",controllers.Login)
}
