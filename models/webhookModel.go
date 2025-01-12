package models

import (
	"time"
)

type Webhook struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null" validate:"required"`
	Name      string    `json:"name" gorm:"size:255"`
	Token     string    `json:"token" gorm:"size:255"`
	InvokedAt time.Time `json:"invoked_at"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
