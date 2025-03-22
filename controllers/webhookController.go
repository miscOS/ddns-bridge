package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/helpers"
	"github.com/miscOS/ddns-bridge/models"
)

func CreateWebhook(c *gin.Context) {

	user, err := fetchUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var webhook models.Webhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	webhook.UserID = user.ID
	webhook.Token = user.Username + "-" + helpers.RandomString(32)

	if err := helpers.GetValidate().Struct(webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := db.GetDB().Create(&webhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "webhook created",
		"data":    webhook,
	})
}

func DeleteWebhook(c *gin.Context) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return
	}

	if err := db.GetDB().Delete(&webhook, webhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "webhook deleted",
		"data":    webhook,
	})
}

func GetWebhook(c *gin.Context) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhook,
	})
}

func GetWebhooks(c *gin.Context) {

	webhooks, err := fetchWebhooksByContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhooks,
	})
}

func fetchWebhookByContext(c *gin.Context) (webhook models.Webhook, err error) {

	user, err := fetchUserFromContext(c)
	if err != nil {
		return webhook, err
	}

	webhookID, err := strconv.ParseUint(c.Param("webhook_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return webhook, err
	}

	webhook.ID = uint(webhookID)
	webhook.UserID = user.ID

	if err := db.GetDB().First(&webhook, webhook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return webhook, err
	}

	return webhook, nil
}

func fetchWebhooksByContext(c *gin.Context) (webhooks []models.Webhook, err error) {

	user, err := fetchUserFromContext(c)
	if err != nil {
		return webhooks, err
	}

	if err := db.GetDB().Where("user_id = ?", user.ID).Find(&webhooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return webhooks, err
	}

	return webhooks, nil
}
