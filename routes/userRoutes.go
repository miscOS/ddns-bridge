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

func ApiRoutes(r *gin.Engine) {

	r.GET("/update", controllers.Update)

	// Secure user routes
	rg := r.Group("/api")
	{
		rg.Use(middleware.UserAuthenticate)
		rg.GET("/user", controllers.GetUser)

		// Webhook routes
		rg.GET("/webhook", controllers.GetWebhooks)
		rg.GET("/webhook/:wid", controllers.GetWebhook)
		rg.POST("/webhook", controllers.CreateWebhook)
		rg.DELETE("/webhook/:wid", controllers.DeleteWebhook)

		// Provider routes
		rg.GET("/webhook/:wid/provider", controllers.GetProviders)
		rg.GET("/webhook/:wid/provider/:pid", controllers.GetProvider)
		rg.POST("/webhook/:wid/provider", controllers.CreateProvider)
		rg.DELETE("/webhook/:wid/provider/:pid", controllers.DeleteProvider)
	}
}
