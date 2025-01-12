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
		rg.GET("/webhook/:webhook_id", controllers.GetWebhook)
		rg.POST("/webhook", controllers.CreateWebhook)
		rg.DELETE("/webhook/:webhook_id", controllers.DeleteWebhook)

		// Task routes
		rg.GET("/webhook/:webhook_id/tasks", controllers.GetTasks)
		rg.GET("/webhook/:webhook_id/task/:task_id", controllers.GetTask)
		rg.POST("/webhook/:webhook_id/task", controllers.CreateTask)
		rg.DELETE("/webhook/:webhook_id/task/:task_id", controllers.DeleteTask)
	}
}
