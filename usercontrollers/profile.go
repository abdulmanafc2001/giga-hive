package usercontrollers

import (
	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

type userProfile struct {
	First_Name string `json:"firstname" gorm:"not null" validate:"min=4,max=20"`
	Last_Name  string `json:"lastname" gorm:"not null" validate:"min=4,max=20"`
	User_Name  string `json:"username" gorm:"not null" validate:"min=4,max=20"`
	Email      string `json:"email" gorm:"not null;unique" validate:"email"`
	Phone      string `json:"phone" gorm:"not null" validate:"min=10,max=10"`
}

func UserProfile(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var user userProfile
	if err := database.DB.Table("users").Select("first_name,last_name,user_name,email,phone").
		Where("id = ?", id).Scan(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find data",
		})
		return
	}

	c.JSON(200, gin.H{
		"profile": user,
	})
}

func ChangePassword(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var input models.CPassword
	if err := c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get body",
		})
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(400, gin.H{
			"errror": "Failed to get user data",
		})
		return
	}

	if err := helpers.CheckPassword(user.Password, input.OldPassword); err != nil {
		c.JSON(400, gin.H{
			"error": "Incorrect old password",
		})
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		c.JSON(400, gin.H{
			"error": "Incorrect confirm password",
		})
		return
	}

	pswd, err := helpers.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Update("password", string(pswd)).Error; err != nil {
		c.JSON(400,gin.H{
			"error" :"Failed to update password",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": "successfully updated password",
	})

}
