package models

import (
	auth_model "hr-system-go/internal/auth/models"
	base_model "hr-system-go/internal/base/models"

	department_model "hr-system-go/internal/department/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type User struct {
	base_model.BaseModel
	Name            string    `gorm:"not null"`
	Email           string    `gorm:"index:idx_email_status,unique"`
	Age             int       `gorm:"not null,check:age > 0"`
	Status          string    `gorm:"index:idx_email_status;default:'active'"`
	PasswordEncrypt string    `gorm:"not null"`
	JoinDate        time.Time `gorm:"type:timestamp;default:current_timestamp()"`
	Salary          *float64
	// Relations
	RoleID       *uint
	Role         *auth_model.Role `gorm:"foreignKey:RoleID"`
	DepartmentID *uint
	Department   *department_model.Department `gorm:"foreignKey:DepartmentID"`
}

func ValidScope(db *gorm.DB) *gorm.DB {
	return db.Model(&User{}).Where("status != ?", "removed")
}

func (u *User) GenerateRandomPassword() {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, 12)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.PasswordEncrypt = string(hashedPassword)
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.DepartmentID == nil {
		return nil
	}
	return tx.Model(&department_model.Department{}).Where("id = ?", *u.DepartmentID).Update("employ_count", gorm.Expr("employ_count + ?", 1)).Error
}
