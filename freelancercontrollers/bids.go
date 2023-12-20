package freelancercontrollers

import (
	"time"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

// ShowAllBids godoc
// @Summary Get all active bids
// @Description Get a list of all active bids with end day greater than or equal to today
// @Tags Bids
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {json} SuccessfulResponse bids
// @Failure 400 {json} ErrorResponse "Failed to find all datas"
// @Router /freelancer/bid/showallbid [get]
func ShowAllBids(c *gin.Context) {
	var bids []models.Bid
	current := time.Now()
	today := current.Format("2006-01-02")
	if err := database.DB.Where("end_day >= ?", today).Find(&bids).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find all datas",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       bids,
	}
	helpers.ResponseResult(c, resp)
}

func AuctionForBid(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")
	id := frlncr.(models.Freelancer).Id

	var auction models.Auction
	if err := c.Bind(&auction); err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	var bid models.Bid
	if err := database.DB.Find(&bid, auction.BidId).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed find this bid",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := database.DB.Create(&models.Auction{
		BidId:         auction.BidId,
		AuctionAmount: auction.AuctionAmount,
		FreelancerId:  id,
	}).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to create auction for bid",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	// c.JSON(200, gin.H{
	// 	"success": "successfull auctioned against bid",
	// })
	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully auctioned against bid",
	}
	helpers.ResponseResult(c, resp)

}