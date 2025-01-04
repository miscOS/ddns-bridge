package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/helpers"
	"github.com/miscOS/ddns-bridge/models"
)

func Signup(c *gin.Context) {

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.GetValidate().Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if password, err := helpers.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		user.Password = password
	}

	result := db.GetDB().Create(&user)

	if result.Error == nil {
		c.JSON(http.StatusCreated, gin.H{"success": "user created"})
	} else if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		c.JSON(http.StatusConflict, gin.H{"failed": "username already exists"})
	} else {
		c.JSON(http.StatusConflict, gin.H{"failed": "unknown error"})
	}
}

func Login(c *gin.Context) {

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.GetValidate().Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var dbUser models.User
	if err := db.GetDB().Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check if the password is correct
	if err := helpers.VerifyPassword(user.Password, dbUser.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := helpers.CreateToken(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

	//c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true) // <- Check what to do here
	//c.Redirect(http.StatusSeeOther, "/")
}

func UserData(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user id"})
		return
	}

	var dbUser models.User
	if err := db.GetDB().First(&dbUser, uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": dbUser.ID, "username": dbUser.Username, "email": dbUser.Email, "created_at": dbUser.CreatedAt, "updated_at": dbUser.UpdatedAt})
}
