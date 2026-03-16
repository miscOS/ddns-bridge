package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
)

type WebhookController struct{}

func (c *WebhookController) Create(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	var request models.CreateWebhookRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "BindJSON", err)
		return
	}

	tasks := []models.Task{}
	if request.Tasks != nil {
		tasks = *request.Tasks
	}

	webhook := models.Webhook{
		Owner: user.ID,
		Name:  *request.Name,
		Tasks: tasks,
		Token: uuid.New().String(),
	}

	if err := db.GetDB().Create(&webhook).Error; err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "DB", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhook,
	})
}

func (c *WebhookController) Read(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Webhook ID", err)
		return
	}

	webhook, err := getWebhook(user.ID, uint(id))
	if err != nil {
		errorResponse(ctx, http.StatusNotFound, "Webhook", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhook,
	})
}

func (c *WebhookController) ReadAll(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	webhooks, err := getWebhooks(user.ID)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "Webhooks", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhooks,
	})
}

func (c *WebhookController) Update(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	var request models.UpdateWebhookRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "BindJSON", err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Webhook ID", err)
		return
	}

	webhook, err := getWebhook(user.ID, uint(id))
	if err != nil {
		errorResponse(ctx, http.StatusNotFound, "Webhook", err)
		return
	}

	webhook.Update(request)

	if err := db.GetDB().Save(&webhook).Error; err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "DB", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhook,
	})
}

func (c *WebhookController) Delete(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Webhook ID", err)
		return
	}

	webhook, err := getWebhook(user.ID, uint(id))
	if err != nil {
		errorResponse(ctx, http.StatusNotFound, "Webhook", err)
		return
	}

	if err := db.GetDB().Delete(&webhook, webhook).Error; err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "DB", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    webhook,
	})
}
