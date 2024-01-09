package routes

import (
	"github.com/abdulmanafc2001/gigahive/handlers"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {

	// Routes for handler functions
	r.GET("/", handlers.Home)
	r.GET("/send", handlers.SendAlertToConnectedUsers)
	r.GET("/ws", handlers.WsEndPoint)

	r.Static("/static", "./static")

}
