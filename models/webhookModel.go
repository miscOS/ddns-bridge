package models

import (
	"time"
)

type Webhook struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null" validate:"required"`
	Name      string    `json:"name" gorm:"size:255"`
	Token     string    `json:"token" gorm:"size:255"`
	IPv4      string    `json:"ipv4" gorm:"size:255"`
	IPv6      string    `json:"ipv6" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
