package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	db "github.com/miscOS/ddns-bridge/database"
	routes "github.com/miscOS/ddns-bridge/routes"
)

func main() {

	// Load the environment variables
	godotenv.Load()

	// Initialize the database
	db.Init()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	routes.Routes(router)
	routes.ApiRoutes(router)

	router.Run()
}
