package freelancercontrollers

import (
	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

func ChangeAcceptedAuctionStatus(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")
	frlncrId := frlncr.(models.Freelancer).Id
	var input struct {
		Id     int    `json:"id"`
		Status string `json:"status"`
	}

	if err := c.Bind(&input); err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var auctions models.AcceptedAuction
	if err := database.DB.Where("freelancer_id = ?", frlncrId).First(&auctions).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "your not access this auctions",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := database.DB.Model(&models.AcceptedAuction{}).Where("id = ?", input.Id).Update("status", input.Status).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "your not access this auctions",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully change status",
	}
	helpers.ResponseResult(c, resp)
}

type BookingDetail struct {
	Id            int    `json:"id"`
	Auction_Id    int    `json:"auctionid"`
	User_Name     string `json:"user_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Amount        int    `json:"amount"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentstatus"`
}

func GetBookingDetail(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")
	frlncrId := frlncr.(models.Freelancer).Id

	var bookings []BookingDetail

	if err := database.DB.Table("accepted_auctions").
		Select("accepted_auctions.id,accepted_auctions.auction_id,users.user_name,users.email,users.phone,accepted_auctions.amount,accepted_auctions.status,accepted_auctions.payment_status").
		Joins("INNER JOIN users ON users.id=accepted_auctions.user_id").Where("accepted_auctions.freelancer_id = ?", frlncrId).Scan(&bookings).Error; err != nil {
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
