package models

type Bid struct {
	Id           int    `json:"id" gorm:"primaryKey"`
	Description  string `json:"description"`
	About        string `json:"about"`
	PriceRange   string `json:"pricerange"`
	ExpectedDays string `json:"expecteddays"`
	User_Id      int    `json:"userid"`
	EndDay       string `json:"endday"`
}
