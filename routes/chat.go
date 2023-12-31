package routes

import (
	"context"

	c "github.com/abdulmanafc2001/gigahive/chatcontrollers"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {
	// http.Handle("/", http.FileServer(http.Dir("./frontent")))
	// http.HandleFunc("/ws", manager.serveWS)
	// http.HandleFunc("/login", manager.loginHandler)

	ctx := context.Background()

	manager := c.NewManager(ctx)
	r.LoadHTMLGlob("./templates/*.html")
	r.GET("/", c.LoadHtml)
	r.GET("/ws",manager.ServeWS)
	r.POST("/login",manager.LoginHandler)

}
