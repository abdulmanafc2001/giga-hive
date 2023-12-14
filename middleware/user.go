package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuthentication(c *gin.Context) {
	tokenString, err := c.Cookie("user_token")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to get token string",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SUPER_SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": "error occurse while token parsing",
		})
		c.AbortWithStatus(401)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user models.User
		database.DB.First(&user, claims["sub"])

		if user.IsBlocked {
			c.AbortWithStatus(401)
		}
		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(401)
	}
}
