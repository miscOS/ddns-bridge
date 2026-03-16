package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;size:255" binding:"required"`
	Password  string    `json:"password,omitempty" gorm:"size:255" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u User) MarshalJSON() ([]byte, error) {

	// Create a new struct with the same fields as User
	type user User
	// Remove the password field
	copy := user(u)
	copy.Password = ""

	return json.Marshal(&copy)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
