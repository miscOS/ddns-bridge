package models

import (
	"time"
)

type Task struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	WebhookID         uint      `json:"webhook_id" gorm:"not null" validate:"required"`
	Name              string    `json:"name" gorm:"size:255"`
	Service           string    `json:"service" gorm:"size:255" validate:"required"`
	ServiceParameters string    `json:"service_parameters" gorm:"default:{}"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
