package models

type User struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	First_Name string `json:"firstname" gorm:"not null" validate:"min=4,max=20"`
	Last_Name  string `json:"lastname" gorm:"not null" validate:"min=4,max=20"`
	User_Name  string `json:"username" gorm:"not null" validate:"min=4,max=20"`
	Email      string `json:"email" gorm:"not null;unique" validate:"email"`
	Phone      string `json:"phone" gorm:"not null" validate:"min=10,max=10"`
	Password   string `json:"password" gorm:"not null"`
	IsBlocked  bool   `json:"isblocked" gorm:"default:false"`
	Otp        int    `json:"otp" gorm:"not null"`
	Validate   bool   `json:"validate" gorm:"default:false"`
}
type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserList struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	First_Name string `json:"firstname" gorm:"not null" validate:"min=4,max=20"`
	Last_Name  string `json:"lastname" gorm:"not null" validate:"min=4,max=20"`
	User_Name  string `json:"username" gorm:"not null" validate:"min=4,max=20"`
	Email      string `json:"email" gorm:"not null;unique" validate:"email"`
	Phone      string `json:"phone" gorm:"not null" validate:"min=10,max=10"`
	IsBlocked  bool   `json:"isblocked" gorm:"default:false"`
	Validate   bool   `json:"validate" gorm:"default:false"`
}

type CPassword struct {
	OldPassword     string `json:"oldpassword"`
	NewPassword     string `json:"newpassword"`
	ConfirmPassword string `json:"confirmpassword"`
}
