package freelancercontrollers

import (
	"time"

	"github.com/abdulmanafc2001/gigahive/database"
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
		c.JSON(400, gin.H{
			"error": "Failed to find all datas",
		})
		return
	}
	c.JSON(200, gin.H{
		"bids": bids,
	})
}


func AuctionForBid(c *gin.Context) {
	frlncr, _ := c.Get("freelancer")
	id := frlncr.(models.Freelancer).Id

	var auction models.Auction
	if err := c.Bind(&auction); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get body",
		})
		return
	}

	var bid models.Bid
	if err := database.DB.Find(&bid,auction.BidId).Error; err != nil {
		c.JSON(400,gin.H{
			"error":"Failed to find this bid",
		})
		return
	}

	if err := database.DB.Create(&models.Auction{
		BidId:         auction.BidId,
		AuctionAmount: auction.AuctionAmount,
		FreelancerId:  id,
	}).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to create auction for bid",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": "successfull auctioned against bid",
	})

}
