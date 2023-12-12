package routes

import (
	controllers "github.com/abdulmanafc2001/gigahive/freelancercontrollers"
	"github.com/gin-gonic/gin"
)

func FreelancerRoutes(router *gin.Engine) {
	r := router.Group("/freelancer")
	r.POST("/signup",controllers.FreelancerSignup)
}
