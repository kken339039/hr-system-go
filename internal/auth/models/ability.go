package models

import (
	base_model "hr-system-go/internal/base/models"
)

type Ability struct {
	base_model.BaseModel
	Name         string `gorm:"not null"`
	Status       string `gorm:"default:'active'"`
	Descriptions string `gorm:"type:text"`
}
