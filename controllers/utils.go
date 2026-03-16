package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
)

func errorResponse(ctx *gin.Context, statusCode int, message string, err error) {
	ctx.JSON(statusCode, gin.H{
		"success": false,
		"error":   fmt.Sprintf("%s: %v", message, err),
	})
}

func getTask(owner uuid.UUID, id uint) (task *models.Task, err error) {
	task = &models.Task{ID: id}

	if err := db.GetDB().Preload("Webhook").First(&task, task).Error; err != nil {
		return nil, err
	}
	if task.Webhook.Owner != owner {
		return nil, fmt.Errorf("you do not have access to this task")
	}
	return task, nil
}

func getTasksByWebhook(webhook *models.Webhook) (tasks *[]models.Task, err error) {

	if err := db.GetDB().Where("webhook_id = ?", webhook.ID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func getUser(ctx *gin.Context) (user *models.User, err error) {

	userContext, ok := ctx.Get("user")
	if !ok {
		return nil, fmt.Errorf("could not retrieve user from context")
	}
	user = userContext.(*models.User)
	return user, nil
}

func getWebhook(owner uuid.UUID, id uint) (webhook *models.Webhook, err error) {

	webhook = &models.Webhook{
		ID:    id,
		Owner: owner,
	}

	if err := db.GetDB().Preload("Tasks").First(&webhook, webhook).Error; err != nil {
		return nil, err
	}
	return webhook, nil
}

func getWebhooks(owner uuid.UUID) (webhooks *[]models.Webhook, err error) {

	if err := db.GetDB().Where("user_id = ?", owner.String()).Find(&webhooks).Error; err != nil {
		return nil, err
	}

	return webhooks, nil
}
