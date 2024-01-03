package routes

import (
	m "github.com/abdulmanafc2001/gigahive/middleware"
	controllers "github.com/abdulmanafc2001/gigahive/usercontrollers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	r := router.Group("/user")
	{
		r.POST("/signup", controllers.UserSignup)
		r.POST("/signup/otpverification", controllers.OtpVerification)
		r.POST("/login", controllers.Login)
		r.POST("/logout", m.UserAuthentication, controllers.Logout)

		r.GET("/profile", m.UserAuthentication, controllers.UserProfile)
		r.PUT("/profile/changepassword", m.UserAuthentication, controllers.ChangePassword)

		r.POST("/bid/createbid", m.UserAuthentication, controllers.CreateBid)
		r.GET("/bid/auctionedbid", m.UserAuthentication, controllers.GetAuctionedBid)

		r.POST("/bid/acceptauction/:auction_id", m.UserAuthentication, controllers.AcceptingEffectiveBid)

		r.GET("/bookingdetails", m.UserAuthentication, controllers.GetBookingDetail)
	}
}
