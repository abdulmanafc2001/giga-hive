package models

type Bid struct {
	Id           int    `json:"id" gorm:"primaryKey"`
	Description  string `json:"description"`
	About        string `json:"about"`
	MinPrice     int    `json:"minprice"`
	MaxPrice     int    `json:"maxprice"`
	ExpectedDays string `json:"expecteddays"`
	User_Id      int    `json:"userid"`
	EndDay       string `json:"endday"`
}

type Auction struct {
	Id            int    `json:"id" gorm:"primaryKey"`
	BidId         int    `json:"bidid"`
	AuctionAmount int `json:"auctionamount"`
	FreelancerId  int    `json:"freelancerid"`
}
