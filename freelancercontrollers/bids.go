package freelancercontrollers

import (
	"time"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

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
	
}
