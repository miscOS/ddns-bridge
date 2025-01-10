package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;size:255" validate:"required"`
	Password  string    `json:"password,omitempty" gorm:"size:255" validate:"required"`
	Email     string    `json:"email" gorm:"size:255" validate:"omitempty,email"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u User) MarshalJSON() ([]byte, error) {

	// Create a new struct with the same fields as User
	type user User
	// Remove the password field
	tmp := user(u)
	tmp.Password = ""

	return json.Marshal(&tmp)
}
