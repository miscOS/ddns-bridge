package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;size:255" validate:"required"`
	Password  string    `json:"password" gorm:"size:255" validate:"required"`
	Email     string    `json:"email" gorm:"size:255" validate:"omitempty,email"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
