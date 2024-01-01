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
