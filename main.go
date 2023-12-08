package main

import (
	"log"

	"github.com/abdulmanafc2001/gigahive/database"
	_ "github.com/abdulmanafc2001/gigahive/docs"
	"github.com/abdulmanafc2001/gigahive/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GigaHive
// @version 1.0
// @description Freelance application API in go using Gin frame work

// @host 	localhost:7000
// @BasePath /
func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}
	database.ConnectToDB()
}

func main() {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run(":7000")
}
