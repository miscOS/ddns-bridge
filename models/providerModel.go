package models

import (
	"time"
)

type Provider struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	WebhookID uint      `json:"webhook_id" gorm:"not null" validate:"required"`
	Name      string    `json:"name" gorm:"size:255"`
	Provider  string    `json:"provider" gorm:"size:255" validate:"required"`
	Settings  string    `json:"settings"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
