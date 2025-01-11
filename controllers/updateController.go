package controllers

import (
	"log"
	"net/http"
	"net/netip"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
	DNSProviders "github.com/miscOS/ddns-bridge/providers"
)

func Update(c *gin.Context) {

	// Parse the IP addresses
	var ipv4 netip.Addr
	var ipv6 netip.Addr
	var err error

	if c.Query("ipv4") != "" {
		ipv4, err = netip.ParseAddr(c.Query("ipv4"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	if c.Query("ipv6") != "" {
		ipv6, err = netip.ParseAddr(c.Query("ipv6"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	params := &DNSProviders.DNSParams{IPv4: ipv4, IPv6: ipv6}

	// Find the webhook
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
		s.Update(params)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "webhook invoked",
	})
}
