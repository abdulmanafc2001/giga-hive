package usercontrollers

import (
	"strconv"
	"time"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

func CreateBid(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var input bid
	if err := c.Bind(&input); err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	endingDate := time.Now().Add(time.Hour * 24 * time.Duration(input.EndDay))
	endDate := endingDate.Format("2006-01-02")

	if err := database.DB.Create(&models.Bid{
		Description:  input.Description,
		About:        input.About,
		MinPrice:     input.MinPrice,
		MaxPrice:     input.MaxPrice,
		ExpectedDays: input.ExpectedDays,
		User_Id:      id,
		EndDay:       endDate,
	}).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to create bid",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully created new bid",
	}
	helpers.ResponseResult(c, resp)
}

type AuctionDetail struct {
	Id            int    `json:"id" gorm:"primaryKey"`
	BidId         int    `json:"bidid"`
	AuctionAmount int    `json:"auctionamount"`
	Full_Name     string `json:"fullname"`
	Description   string `json:"description"`
	About         string `json:"about"`
	MinPrice      int    `json:"minprice"`
	MaxPrice      int    `json:"maxprice"`
	ExpectedDays  string `json:"expecteddays"`
	EndDay        string `json:"endday"`
}

func GetAuctionedBid(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var auctions []AuctionDetail
	if err := database.DB.Table("auctions").Select("auctions.id,auctions.bid_id,auctions.auction_amount,freelancers.full_name,bids.description,bids.about,bids.min_price,bids.max_price,bids.expected_days,bids.end_day").
		Joins("INNER JOIN freelancers ON freelancers.id=auctions.freelancer_id").
		Joins("INNER JOIN bids ON bids.id=auctions.bid_id").Where("bids.user_id=?", id).Scan(&auctions).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find auctioned datas",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       auctions,
	}
	helpers.ResponseResult(c, resp)
}

// user is accepting effective auction against the bid
func AcceptingEffectiveBid(c *gin.Context) {
	usr, _ := c.Get("user")
	usrId := usr.(models.User).Id

	id, _ := strconv.Atoi(c.Param("auction_id"))

	var auction models.Auction
	if err := database.DB.Where("user_id = ? AND id = ?", usrId, id).First(&auction).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find this auction",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	database.DB.Create(&models.AcceptedAuction{
		User_Id: usrId,
		Freelancer_Id: auction.FreelancerId,
		Amount: auction.AuctionAmount,
		Status: "Pending",
	})

}
