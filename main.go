package main

import (
	"log"

	"11val/config"
	"11val/models"
	"11val/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	models.DB = config.ConnectDatabase()

	models.DB.AutoMigrate(&models.User{})

	router := gin.Default()

	routes.UserRoutes(router)

	log.Println("Server running on port 8080")

	router.Run(":8080")
}
