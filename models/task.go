package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID              uint                   `json:"id" gorm:"primaryKey"`
	Webhook         *Webhook               `json:"webhook" gorm:"foreignKey:WebhookID;constraint:OnDelete:CASCADE;"`
	WebhookID       uint                   `json:"-"`
	Name            string                 `json:"name" gorm:"size:255"`
	Service         string                 `json:"service" gorm:"size:255"`
	ServiceParams   map[string]interface{} `json:"service_params" gorm:"-"`
	ServiceParamsDB string                 `json:"-"` // Database field
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateTaskRequest struct {
	Name          *string                 `json:"name" binding:"required"`
	WebhookID     *uint                   `json:"webhook_id" binding:"required"`
	Service       *string                 `json:"service" binding:"required"`
	ServiceParams *map[string]interface{} `json:"service_params" binding:"-"`
}

type UpdateTaskRequest struct {
	Name          *string                 `json:"name" binding:"-"`
	Service       *string                 `json:"service" binding:"-"`
	ServiceParams *map[string]interface{} `json:"service_params" binding:"-"`
}

func (t *Task) BeforeSave(tx *gorm.DB) (err error) {
	bytes, err := json.Marshal(t.ServiceParams)
	if err != nil {
		return err
	}
	t.ServiceParamsDB = string(bytes)
	return nil
}

func (t *Task) AfterFind(tx *gorm.DB) (err error) {
	err = json.Unmarshal([]byte(t.ServiceParamsDB), &t.ServiceParams)
	return err
}

func (t *Task) Update(request UpdateTaskRequest) {

	if request.Name != nil {
		t.Name = *request.Name
	}
	if request.Service != nil {
		t.Service = *request.Service
	}
	if request.ServiceParams != nil {
		t.ServiceParams = *request.ServiceParams
	}
}
