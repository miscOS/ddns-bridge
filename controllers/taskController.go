package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/miscOS/ddns-bridge/database"
	"github.com/miscOS/ddns-bridge/helpers"
	"github.com/miscOS/ddns-bridge/models"
)

func CreateTask(c *gin.Context) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return
	}

	task := models.Task{}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	task.WebhookID = webhook.ID

	if err := helpers.GetValidate().Struct(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := db.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "task created",
		"data":    task,
	})
}

func DeleteTask(c *gin.Context) {

	task, err := fetchTaskByContext(c)
	if err != nil {
		return
	}

	if err := db.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "task deleted",
		"data":    task,
	})
}

func GetTask(c *gin.Context) {

	task, err := fetchTaskByContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

func GetTasks(c *gin.Context) {

	tasks, err := fetchTasksByContext(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tasks,
	})
}

func fetchTaskByContext(c *gin.Context) (task models.Task, err error) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return task, err
	}

	taskID, err := strconv.ParseUint(c.Param("task_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return task, err
	}

	task.ID = uint(taskID)
	task.WebhookID = webhook.ID

	if err := db.GetDB().First(&task, task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return task, err
	}

	return task, nil
}

func fetchTasksByContext(c *gin.Context) (tasks []models.Task, err error) {

	webhook, err := fetchWebhookByContext(c)
	if err != nil {
		return tasks, err
	}

	if err := db.GetDB().Where("webhook_id = ?", webhook.ID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return tasks, err
	}
	return tasks, nil
}

func fetchTasksByWebhook(webhook *models.Webhook) (tasks []models.Task, err error) {

	if err := db.GetDB().Where("webhook_id = ?", webhook.ID).Find(&tasks).Error; err != nil {
		return tasks, err
	}
	return tasks, nil
}
