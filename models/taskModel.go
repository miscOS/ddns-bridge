package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID               uint                   `json:"id" gorm:"primaryKey"`
	WebhookID        uint                   `json:"webhook_id" gorm:"not null" validate:"required"`
	Name             string                 `json:"name" gorm:"size:255"`
	Service          string                 `json:"service" gorm:"size:255" validate:"required"`
	ServiceParams    map[string]interface{} `json:"service_params" gorm:"-"`
	ServiceParamsRaw string                 `json:"-"` // Database field
	CreatedAt        time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

func (t *Task) BeforeSave(tx *gorm.DB) (err error) {
	bytes, err := json.Marshal(t.ServiceParams)
	if err != nil {
		return err
	}
	t.ServiceParamsRaw = string(bytes)
	return nil
}

func (t *Task) AfterFind(tx *gorm.DB) (err error) {
	err = json.Unmarshal([]byte(t.ServiceParamsRaw), &t.ServiceParams)
	return err
}
