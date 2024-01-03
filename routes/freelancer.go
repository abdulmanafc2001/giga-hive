package routes

import (
	controllers "github.com/abdulmanafc2001/gigahive/freelancercontrollers"
	m "github.com/abdulmanafc2001/gigahive/middleware"
	"github.com/gin-gonic/gin"
)

func FreelancerRoutes(router *gin.Engine) {
	r := router.Group("/freelancer")
	{
		r.POST("/signup", controllers.FreelancerSignup)
		r.POST("/signup/otpverification", controllers.ValidateOTP)
		r.POST("/login", controllers.Login)
		r.POST("/logout", m.FreelancerAuthentication, controllers.Logout)

		r.GET("/profile", m.FreelancerAuthentication, controllers.GetProfile)
		r.PUT("/profile/changepassword", m.FreelancerAuthentication, controllers.ChangePassword)

		r.GET("/bid/showallbid", m.FreelancerAuthentication, controllers.ShowAllBids)
		r.POST("/bid/auction", m.FreelancerAuthentication, controllers.AuctionForBid)
		r.PUT("/auction/changestatus", m.FreelancerAuthentication, controllers.ChangeAcceptedAuctionStatus)

		r.GET("/bookingdetails", m.FreelancerAuthentication, controllers.GetBookingDetail)
	}

}
