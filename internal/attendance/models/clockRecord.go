package models

import (
	user_models "hr-system-go/internal/user/models"
	"time"
)

type ClockRecord struct {
	UserID   uint
	User     user_models.User `gorm:"foreignKey:UserID"`
	ClockIn  time.Time
	ClockOut time.Time
	Date     time.Time
	Status   string `gorm:"default:'present'"`
	Notes    string `gorm:"type:varchar(255)"`
}
