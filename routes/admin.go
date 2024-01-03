package routes

import (
	controllers "github.com/abdulmanafc2001/gigahive/admincontrollers"
	m "github.com/abdulmanafc2001/gigahive/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	r := router.Group("/admin")
	r.POST("/login", controllers.Login)
	r.POST("/logout", m.AdminAuthentication, controllers.Logout)
	r.GET("/user/list", m.AdminAuthentication, controllers.GetUserList)
	r.GET("/freelancer/list", m.AdminAuthentication, controllers.ListFreelancers)
	r.PATCH("/user/block", m.AdminAuthentication, controllers.BlockUser)
	r.PATCH("/user/unblock", m.AdminAuthentication, controllers.UnBlockUser)
}
