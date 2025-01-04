package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/miscOS/ddns-bridge/controllers"
	"github.com/miscOS/ddns-bridge/middleware"
)

func UserRoutes(r *gin.Engine) {
	r.POST("users/signup", controller.Signup)
	r.POST("users/login", controller.Login)
}

func SecureUserRoutes(r *gin.Engine) {
	r.Use(middleware.UserAuthenticate)
	r.GET("/users/data", controller.UserData)
	//r.GET("/users/:user_id", controller.GetUser())
}
