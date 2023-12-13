package routes

import (
	controllers "github.com/abdulmanafc2001/gigahive/freelancercontrollers"
	"github.com/abdulmanafc2001/gigahive/middleware"
	"github.com/gin-gonic/gin"
)

func FreelancerRoutes(router *gin.Engine) {
	r := router.Group("/freelancer")
	r.POST("/signup", controllers.FreelancerSignup)
	r.POST("/signup/otpverification", controllers.ValidateOTP)
	r.POST("/login", controllers.Login)
	r.GET("/profile", middleware.FreelancerAuthentication, controllers.GetProfile)
}
