package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	StatusCode int `json:"statuscode"`
	Err        any `json:"error"`
	Data       any `json:"data"`

}

func ResponseResult(c *gin.Context,response Response) {

	c.JSON(response.StatusCode,response)
}
