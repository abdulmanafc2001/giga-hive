package main

import (
	"log"

	"github.com/abdulmanafc2001/gigahive/database"
	"github.com/abdulmanafc2001/gigahive/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}
	database.ConnectToDB()
}

func main() {
	router := gin.Default()
	routes.UserRoutes(router)
	router.Run(":7000")
}
