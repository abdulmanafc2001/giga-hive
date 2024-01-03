package usercontrollers

import (
	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

type BookingDetail struct {
	Id            int    `json:"id"`
	Auction_Id    int    `json:"auctionid"`
	Full_Name     string `json:"freelancer_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Amount        int    `json:"amount"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentstatus"`
}

func GetBookingDetail(c *gin.Context) {
	usr, _ := c.Get("user")
	usrId := usr.(models.User).Id

	var bookings []BookingDetail

	if err := database.DB.Table("accepted_auctions").
		Select("accepted_auctions.id,accepted_auctions.auction_id,freelancers.full_name,freelancers.email,freelancers.phone,accepted_auctions.amount,accepted_auctions.status,accepted_auctions.payment_status").
		Joins("INNER JOIN freelancers ON freelancers.id=accepted_auctions.freelancer_id").Where("accepted_auctions.user_id = ?", usrId).Scan(&bookings).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "failed to get booking details",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       bookings,
	}
	helpers.ResponseResult(c, resp)
}
