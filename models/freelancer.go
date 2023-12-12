package models

type Freelancer struct {
	Id            int    `json:"id" gorm:"primaryKey"`
	Full_Name     string `json:"fullname" validate:"min=4,max=20"`
	User_Name     string `json:"username" validate:"min=4,max=20"`
	Email         string `json:"email" validate:"email"`
	Password      string `json:"password" validate:"min=4,max=20"`
	Qualification string `json:"qualification"`
	Tools         string `json:"tools"`
	OTP           string `json:"otp"`
	IsBlocked     bool    `json:"isblocked" gorm:"default:false"`
	Validate      bool   `json:"validate" gorm:"default:false"`
}


