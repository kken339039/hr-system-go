package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp()"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp()"`
}

func (model *BaseModel) BeforeCreate(db *gorm.DB) (err error) {
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return
}

func (model *BaseModel) BeforeUpdate(db *gorm.DB) (err error) {
	model.UpdatedAt = time.Now()
	return
}
