package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/models"
)

type TaskController struct{}

func (c *TaskController) Create(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	var request models.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "BindJSON", err)
		return
	}

	webhook, err := getWebhook(user.ID, *request.WebhookID)
	if err != nil {
		errorResponse(ctx, http.StatusNotFound, "Webhook", err)
		return
	}

	task := models.Task{
		WebhookID:     webhook.ID,
		Name:          *request.Name,
		Service:       *request.Service,
		ServiceParams: *request.ServiceParams,
	}

	if err := db.GetDB().Create(&task).Error; err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "DB", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

func (c *TaskController) Read(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Task ID", err)
		return
	}

	task := &models.Task{ID: uint(id)}

	if err := db.GetDB().Preload("Webhook").First(&task, task).Error; err != nil {
		errorResponse(ctx, http.StatusNotFound, "DB", err)
		return
	}

	if task.Webhook.Owner != user.ID {
		errorResponse(ctx, http.StatusForbidden, "Access", fmt.Errorf("you do not have access to this task"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

func (c *TaskController) Update(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	var request models.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "BindJSON", err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Task ID", err)
		return
	}

	task, err := getTask(user.ID, uint(id))
	if err != nil {
		errorResponse(ctx, http.StatusNotFound, "Task", err)
		return
	}

	task.Update(request)

	if err := db.GetDB().Save(&task).Error; err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "DB", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

func (c *TaskController) Delete(ctx *gin.Context) {

	user, err := getUser(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "User", err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Task ID", err)
		return
	}

	task, err := getTask(user.ID, uint(id))
	if err != nil {
		errorResponse(ctx, http.StatusNotFound, "Task", err)
		return
	}

	if err := db.GetDB().Delete(&task).Error; err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "DB", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}
