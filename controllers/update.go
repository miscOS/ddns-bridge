package controllers

import (
	"net/http"
	"net/netip"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
	"github.com/miscOS/ddns-bridge/services"
)

func Update(c *gin.Context) {

	values := &models.UpdateRequest{}
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

	tasks, err := getTasksByWebhook(webhook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var result []models.UpdateResult

	for _, task := range *tasks {

		service, err := services.Get(task.Service)
		if err != nil {
			result = append(result, models.UpdateResult{Success: false, Message: "Service not found: " + task.Service})
			continue
		}

		if err := service.Init(task.ServiceParams); err != nil {
			result = append(result, models.UpdateResult{Success: false, Message: "Error initializing service: " + err.Error()})
			continue
		}

		res, err := service.Update(values)
		if err != nil {
			result = append(result, models.UpdateResult{Success: false, Message: "Error updating service: " + err.Error()})
		}

		for _, r := range res {
			r.Task = task.Name
			result = append(result, r)
		}
	}

	c.JSON(http.StatusOK, result)
}
