package main

import (
	"github.com/gin-gonic/gin"

	db "github.com/miscOS/ddns-bridge/database"
	routes "github.com/miscOS/ddns-bridge/routes"
)

func main() {

	// Initialize the database
	db.Init()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	routes.UserRoutes(router)
	routes.SecureUserRoutes(router)

	router.Run()
}
