package seeds

import (
	auth_constants "hr-system-go/internal/auth/constants"
	auth_models "hr-system-go/internal/auth/models"
	department_constants "hr-system-go/internal/department/constants"
	department_models "hr-system-go/internal/department/models"
	"hr-system-go/internal/user/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	Seeds = append(Seeds, Seed{
		Name: "20240711213135-import-hr-user",
		Exec: Exec_20240711213135,
	})
}

func Exec_20240711213135(db *gorm.DB) error {
	// TODO: Implement the seed logic here
	user := &models.User{
		Name:  "HR Manager User",
		Email: "hr-manager@gmail.com",
		Age:   31,
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("HrManager123"), bcrypt.DefaultCost)
	user.PasswordEncrypt = string(hashedPassword)
	user.JoinDate = time.Now()

	var hrManagerRole *auth_models.Role
	if err := db.Where("name = ?", auth_constants.ROLE_HR_MANAGER).First(&hrManagerRole).Error; err != nil {
		return err
	}
	user.RoleID = &hrManagerRole.ID
	var hrDepartment *department_models.Department
	if err := db.Where("name = ?", department_constants.DEPARTMENT_HR).First(&hrDepartment).Error; err != nil {
		return err
	}
	user.DepartmentID = &hrDepartment.ID
	result := db.Where("email = ?", user.Email).FirstOrCreate(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
