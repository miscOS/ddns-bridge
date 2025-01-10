package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/helpers"
	"github.com/miscOS/ddns-bridge/models"
)

func CreateProvider(c *gin.Context) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return
	}

	provider := models.Provider{}
	if err := c.BindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	provider.WebhookID = webhook.ID

	if err := helpers.GetValidate().Struct(provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := db.GetDB().Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "provider created",
		"data":    provider,
	})
}

func DeleteProvider(c *gin.Context) {

	provider, err := fetchProviderByContext(c)
	if err != nil {
		return
	}

	if err := db.GetDB().Delete(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "provider deleted",
	})
}

func GetProvider(c *gin.Context) {

	provider, err := fetchProviderByContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    provider,
	})
}

func GetProviders(c *gin.Context) {

	providers, err := fetchProvidersByContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    providers,
	})
}

func fetchProviderByContext(c *gin.Context) (provider models.Provider, err error) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return provider, err
	}

	providerID, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return provider, err
	}

	provider.ID = uint(providerID)
	provider.WebhookID = webhook.ID

	if err := db.GetDB().First(&provider, provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return provider, err
	}

	return provider, nil
}

func fetchProvidersByContext(c *gin.Context) (providers []models.Provider, err error) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return providers, err
	}

	if err := db.GetDB().Where("webhook_id = ?", webhook.ID).Find(&providers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return providers, err
	}
	return providers, nil
}

func fetchProvidersByWebhook(webhook *models.Webhook) (providers []models.Provider, err error) {

	if err := db.GetDB().Where("webhook_id = ?", webhook.ID).Find(&providers).Error; err != nil {
		return providers, err
	}
	return providers, nil
}
