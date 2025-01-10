package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miscOS/ddns-bridge/controllers"
	"github.com/miscOS/ddns-bridge/middleware"
)

func Routes(r *gin.Engine) {
	// Public user routes
	rg := r.Group("/api/user")
	{
		rg.POST("signup", controllers.Signup)
		rg.POST("login", controllers.Login)
	}
}

func SecureRoutes(r *gin.Engine) {
	// Secure user routes
	rg := r.Group("/api/user")
	{
		rg.Use(middleware.UserAuthenticate)
		rg.GET("/info", controllers.UserInfo)

		// Hook routes
		rg.POST("/hook", controllers.CreateWebhook)
		rg.DELETE("/hook", controllers.DeleteWebhook)
	}
}
