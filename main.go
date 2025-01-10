package main

import (
	"github.com/gin-gonic/gin"

	db "github.com/miscOS/ddns-bridge/database"
	routes "github.com/miscOS/ddns-bridge/routes"
)

func main() {

	// Generate a secret key
	//helpers.GenerateSecret()

	// Initialize the database
	db.Init()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	routes.Routes(router)
	routes.ApiRoutes(router)

	router.Run()
}
