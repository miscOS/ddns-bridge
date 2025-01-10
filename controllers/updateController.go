package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
	DNSProviders "github.com/miscOS/ddns-bridge/providers"
)

func Update(c *gin.Context) {

	webhook := &models.Webhook{Token: c.Query("token")}

	if err := db.GetDB().First(&webhook, webhook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	providers, err := fetchProvidersByWebhook(webhook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	for _, provider := range providers {
		s, err := DNSProviders.GetDNSProvider(provider.Provider)
		if err != nil {
			log.Println(err)
			continue
		}
		s.Setup(provider.Settings)
		s.Update()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "webhook invoked",
	})
}
