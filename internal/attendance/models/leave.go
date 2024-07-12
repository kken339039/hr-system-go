package models

import (
	base_model "hr-system-go/internal/base/models"
	user_model "hr-system-go/internal/user/models"
	"time"

	"gorm.io/gorm"
)

type Leave struct {
	base_model.BaseModel
	UserID    uint
	User      user_model.User `gorm:"foreignKey:UserID"`
	StartDate time.Time       `gorm:"type:timestamp;not null"`
	EndDate   time.Time       `gorm:"type:timestamp;not null"`
	LeaveType string          `gorm:"not null"`
	Status    string          `gorm:"default:'pending'"`
}

func ValidLeaveScope(db *gorm.DB) *gorm.DB {
	return db.Model(&Leave{}).Where("status != ?", "removed")
}
