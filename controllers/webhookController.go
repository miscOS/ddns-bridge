package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/helpers"
	"github.com/miscOS/ddns-bridge/models"
)

func CreateWebhook(c *gin.Context) {

	user, ok := GetUser(c)
	if !ok {
		return
	}

	var webhook models.Webhook
	if err := c.BindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhook.UserID = user.ID
	webhook.Token = user.Username + "-" + helpers.RandomString(32)

	if err := db.GetDB().Create(&webhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "instance created"})
}

func DeleteWebhook(c *gin.Context) {

	user, ok := GetUser(c)
	if !ok {
		return
	}

	var webhook models.Webhook
	if err := c.BindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhook.UserID = user.ID

	if err := db.GetDB().Delete(&webhook, webhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "instance deleted"})
}
