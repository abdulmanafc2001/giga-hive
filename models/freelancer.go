package models

type Freelancer struct {
	Id              int    `json:"id" gorm:"primaryKey"`
	Full_Name       string `json:"fullname" validate:"min=4,max=20"`
	User_Name       string `json:"username" validate:"min=4,max=20"`
	Email           string `json:"email" validate:"email"`
	Phone           string `json:"phone" validate:"min=10,max=10"`
	Password        string `json:"password" validate:"min=4,max=20"`
	Qualification   string `json:"qualification"`
	Tools           string `json:"tools"`
	Rating          string `json:"rating"`
	NumberOfRatings int    `json:"numberofratings"`
	OTP             string `json:"otp"`
	IsBlocked       bool   `json:"isblocked" gorm:"default:false"`
	Validate        bool   `json:"validate" gorm:"default:false"`
}
