package models

import (
	"time"

	"github.com/google/uuid"
)

type Webhook struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:255"`
	Owner     uuid.UUID `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"size:255"`
	IPv4      string    `json:"ipv4" gorm:"size:255"`
	IPv6      string    `json:"ipv6" gorm:"size:255"`
	Tasks     []Task    `json:"tasks" gorm:"foreignKey:WebhookID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateWebhookRequest struct {
	Name  *string `json:"name" binding:"required"`
	Tasks *[]Task `json:"tasks" binding:"-"`
}

type UpdateWebhookRequest struct {
	Name  *string `json:"name" binding:"-"`
	Token *string `json:"token" binding:"-"`
	Tasks *[]Task `json:"tasks" binding:"-"`
}

func (w *Webhook) Update(request UpdateWebhookRequest) {

	if request.Name != nil {
		w.Name = *request.Name
	}

	if request.Token != nil && *request.Token != w.Token {
		w.Token = uuid.New().String()
	}
}
