package usercontrollers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	auctionId, _ := strconv.Atoi(c.Param("auction_id"))
	id := 1
	// var user models.User
	// err := database.DB.First(&user, id).Error
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": "This user didn't find",
	// 	})
	// 	return
	// }

	var auction models.AcceptedAuction
	err := database.DB.Where("id = ? AND user_id = ?", auctionId, id).First(&auction).Error
	if err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Auction not found",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
	}
	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"), os.Getenv("RAZOR_SECRET"))
	data := map[string]interface{}{
		"amount":   auction.Amount * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	value := body["id"]
	c.HTML(http.StatusOK, "app.html", gin.H{
		"userid":       id,
		"freelancerid": auction.Freelancer_Id,
		"totalprice":   auction.Amount,
		"paymentid":    value,
		"auctionid":    auctionId,
	})
}

func SuccessPayment(c *gin.Context) {
	// user, _ := c.Get("user")
	userId, _ := strconv.Atoi(c.Query("user_id"))

	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")
	freelancerId, _ := strconv.Atoi(c.Query("freelancer_id"))
	auctionId, _ := strconv.Atoi(c.Query("auction_id"))

	err := database.DB.Create(&models.RazorPay{
		RazorPayment_id:  paymentid,
		User_id:          userId,
		Freelancer_Id:    freelancerId,
		Signature:        signature,
		RazorPayOrder_id: orderid,
		AmountPaid:       totalamount,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if err := database.DB.Model(&models.AcceptedAuction{}).Where("id = ?", auctionId).Update("payment_status", "Completed").Error; err != nil {
		resp := helpers.Response{
			StatusCode: 500,
			Err:        "Something wrong please try again",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	c.JSON(200, gin.H{
		"status":     true,
		"payment_id": paymentid,
		"Message":    "fskf",
		"notice":     "ksfk",
	})
}

func Success(c *gin.Context) {
	pid := c.Query("id")
	fmt.Println("successs",pid)

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}
