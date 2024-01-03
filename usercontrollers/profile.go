package usercontrollers

import (
	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/helpers"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
)

type userProfile struct {
	Id         int    `json:"id"`
	First_Name string `json:"firstname"`
	Last_Name  string `json:"lastname"`
	User_Name  string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

func UserProfile(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var user userProfile
	if err := database.DB.Table("users").Select("id,first_name,last_name,user_name,email,phone").
		Where("id = ?", id).Scan(&user).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to find data",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       user,
	}
	helpers.ResponseResult(c, resp)
}

func ChangePassword(c *gin.Context) {
	usr, _ := c.Get("user")
	id := usr.(models.User).Id

	var input models.CPassword
	if err := c.Bind(&input); err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get body",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to get user data",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := helpers.CheckPassword(user.Password, input.OldPassword); err != nil {
		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Incorrect old password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if input.NewPassword != input.ConfirmPassword {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Incorrect confirm password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	pswd, err := helpers.HashPassword(input.NewPassword)
	if err != nil {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to hash password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Update("password", string(pswd)).Error; err != nil {

		resp := helpers.Response{
			StatusCode: 400,
			Err:        "Failed to update password",
			Data:       nil,
		}
		helpers.ResponseResult(c, resp)
		return
	}

	resp := helpers.Response{
		StatusCode: 200,
		Err:        nil,
		Data:       "Successfully updated password",
	}
	helpers.ResponseResult(c, resp)

}
