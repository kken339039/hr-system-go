package models

import (
	user_models "hr-system-go/internal/user/models"
	"time"
)

type ClockRecord struct {
	UserID   uint
	User     user_models.User `gorm:"foreignKey:UserID"`
	ClockIn  time.Time        `gorm:"type:timestamp;default:current_timestamp()"`
	ClockOut *time.Time       `gorm:"default:null"`
}
