package routes

import (
	"github.com/abdulmanafc2001/gigahive/middleware"
	controllers "github.com/abdulmanafc2001/gigahive/usercontrollers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	r := router.Group("/user")
	{
		r.POST("/signup", controllers.UserSignup)
		r.POST("/signup/otpverification", controllers.OtpVerification)
		r.POST("/login", controllers.Login)
		r.GET("/profile", middleware.UserAuthentication, controllers.UserProfile)
		r.PUT("/profile/changepassword", middleware.UserAuthentication, controllers.ChangePassword)
		r.POST("/bid/createbid", middleware.UserAuthentication, controllers.CreateBid)
	}

}
