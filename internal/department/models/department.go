package models

import (
	base_model "hr-system-go/internal/base/models"

	"gorm.io/gorm"
)

type Department struct {
	base_model.BaseModel
	Name         string `gorm:"not null"`
	Descriptions string `gorm:"type:text"`
	Status       string `gorm:"default:'active'"`
	EmployCount  int
}

func ValidScope(db *gorm.DB) *gorm.DB {
	return db.Model(&Department{}).Where("status != ?", "removed")
}

func (d *Department) UpdateEmployCount(db *gorm.DB, change int) error {
	// change might be +-1
	return db.Model(d).Update("employ_count", d.EmployCount+change).Error
}
