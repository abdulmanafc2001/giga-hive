package models

type Bid struct {
	Id           int    `json:"id" gorm:"primaryKey"`
	Description  string `json:"description"`
	About        string `json:"about"`
	MinPrice     int    `json:"minprice"`
	MaxPrice     int    `json:"maxprice"`
	ExpectedDays string `json:"expecteddays"`
	User_Id      int    `json:"userid"`
	EndDay       string `json:"last_day"`
	Auctioned    bool   `json:"auctioned" gorm:"default:false"`
}

type Auction struct {
	Id            int `json:"id" gorm:"primaryKey"`
	User_Id       int `json:"userid"`
	BidId         int `json:"bidid"`
	AuctionAmount int `json:"auctionamount"`
	FreelancerId  int `json:"freelancerid"`
}

type AcceptedAuction struct {
	Id            int    `json:"id" gorm:"primaryKey"`
	Auction_Id    int	`json:"auctionid"`
	User_Id       int    `json:"userid"`
	Freelancer_Id int    `json:"freelancerid"`
	Amount        int    `json:"amount"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentstatus"`
}
