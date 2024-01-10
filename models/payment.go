package models

type RazorPay struct {
	RazorPayment_id  string `json:"razorpaymentid" gorm:"primaryKey"`
	Freelancer_Id    int    `json:"freelancer_id"`
	User_id          int    `json:"userid"`
	RazorPayOrder_id string `json:"razorpayorderid"`
	Signature        string `json:"signature"`
	AmountPaid       string `json:"amountpaid"`
}
