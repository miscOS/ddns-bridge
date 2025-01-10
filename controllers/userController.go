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

	token, err := helpers.CreateToken(dbUser.Username, helpers.GetSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

	//c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true) // <- Check what to do here
	//c.Redirect(http.StatusSeeOther, "/")
}

func GetUser(c *gin.Context) {

	user, err := fetchUserFromContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func fetchUserFromContext(c *gin.Context) (user models.User, err error) {
	username := c.GetString("username")

	if err := db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return user, err
	}
	return user, nil
}
