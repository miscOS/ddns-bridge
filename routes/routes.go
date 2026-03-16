package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miscOS/ddns-bridge/auth"
	"github.com/miscOS/ddns-bridge/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	// Healthcheck Route
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "running",
		})
	})

	r.GET("/update", controllers.Update)

	// User Routes (only if local JWT auth is enabled)
	localAuthenticator, err := auth.GetLocalAuthenticator()
	if err == nil {

		user := &controllers.UserController{
			Authenticator: localAuthenticator,
		}

		rg := r.Group("/auth")
		{
			rg.POST("/login", user.Login)
			rg.POST("/register", user.Create)
		}
	}
}

func RegisterSecureRoutes(r *gin.Engine) {

	// User profile route
	user := &controllers.UserController{}
	r.GET("/profile", auth.GetMiddleware(), user.Read)

	// API routes
	rg := r.Group("/api")
	{
		rg.Use(auth.GetMiddleware())

		// Webhook routes
		webhook := &controllers.WebhookController{}
		rg.POST("/webhook", webhook.Create)
		rg.GET("/webhook/:id", webhook.Read)
		rg.GET("/webhook", webhook.ReadAll)
		rg.PATCH("/webhook/:id", webhook.Update)
		rg.DELETE("/webhook/:id", webhook.Delete)

		// Task routes
		task := &controllers.TaskController{}
		rg.POST("/task", task.Create)
		rg.GET("/task/:id", task.Read)
		rg.PATCH("/task/:id", task.Update)
		rg.DELETE("/task/:id", task.Delete)
	}
}
