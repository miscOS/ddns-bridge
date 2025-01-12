package controllers

import (
	"log"
	"net/http"
	"net/netip"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/helpers"
	"github.com/miscOS/ddns-bridge/models"
)

func Update(c *gin.Context) {

	values := &models.UpdaetValue{}
	webhook := &models.Webhook{Token: c.Query("token")}

	if err := db.GetDB().First(&webhook, webhook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if c.Query("ipv4") != "" {
		ipv4, err := netip.ParseAddr(c.Query("ipv4"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		values.IPv4 = ipv4
		webhook.IPv4 = ipv4.String()
	}

	if c.Query("ipv6") != "" {
		ipv6, err := netip.ParseAddr(c.Query("ipv6"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		values.IPv6 = ipv6
		webhook.IPv6 = ipv6.String()
	}

	db.GetDB().Save(&webhook)

	providers, err := fetchTasksByWebhook(webhook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var updateResult []models.UpdateResult

	for _, provider := range providers {

		if s := helpers.GetService(provider.Service); s != nil {

			if err := s.Setup(provider.ServiceParameters); err != nil {
				continue
			}

			res, err := s.Update(values)
			if err != nil {
				log.Printf("Error updating DNS for provider %s: %s", provider.Service, err.Error())
			}
			updateResult = append(updateResult, res...)

		}
	}

	c.JSON(http.StatusOK, updateResult)
}
