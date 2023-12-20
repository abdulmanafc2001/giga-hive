package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func FreelancerAuthentication(c *gin.Context) {
	tokenStr := c.Request.Header.Get("freelancer_token")
	if len(tokenStr) == 0 {
		err := errors.New("autharization header not provided")
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SUPER_SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": "error occurse while token generation",
		})
		c.AbortWithStatus(401)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var freelancer models.Freelancer
		database.DB.First(&freelancer, claims["sub"])

		if freelancer.IsBlocked {
			c.AbortWithStatus(401)
		}
		c.Set("freelancer", freelancer)

		c.Next()
	} else {
		c.AbortWithStatus(401)
	}
}
