package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/miscOS/ddns-bridge/auth"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
)

type UserController struct {
	Authenticator *auth.LocalAuthenticator
}

func (c *UserController) Login(ctx *gin.Context) {

	var request models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user models.User
	if err := db.GetDB().Where("username = ?", request.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check if the password is correct
	if err := c.Authenticator.VerifyPassword(request.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := c.Authenticator.GenerateToken(user.ID.String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) Create(ctx *gin.Context) {

	if ok := c.Authenticator.AllowRegistration(ctx.Query("key")); !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "registration is not allowed"})
		return
	}

	var request models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if password, err := c.Authenticator.HashPassword(request.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		request.Password = password
	}

	user := models.User{
		ID:       uuid.New(),
		Username: request.Username,
		Password: request.Password,
	}

	result := db.GetDB().Create(&user)

	if result.Error == nil {
		ctx.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    user,
		})
	} else {
		ctx.JSON(http.StatusConflict, gin.H{"error": result.Error.Error()})
	}
}

func (c *UserController) Read(ctx *gin.Context) {

	request, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user from context"})
		return
	}

	user := request.(*models.User)

	if err := db.GetDB().First(&user, user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}
