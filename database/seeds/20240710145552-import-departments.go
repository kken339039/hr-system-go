package seeds

import (
	"hr-system-go/internal/department/constants"
	"hr-system-go/internal/department/models"

	"gorm.io/gorm"
)

func init() {
	Seeds = append(Seeds, Seed{
		Name: "20240710145552-import-departments",
		Exec: Exec_20240710145552,
	})
}

func Exec_20240710145552(db *gorm.DB) error {
	departments := []models.Department{
		{
			Name:         constants.DEPARTMENT_ADMIN,
			Descriptions: "Administration",
		},
		{
			Name:         constants.DEPARTMENT_RD,
			Descriptions: "Research and Development",
		},
		{
			Name:         constants.DEPARTMENT_BD,
			Descriptions: "Business Development",
		},
		{
			Name:         constants.DEPARTMENT_MKT,
			Descriptions: "MARKETING",
		},
		{
			Name:         constants.DEPARTMENT_HR,
			Descriptions: "Human Resource",
		},
	}

	for _, department := range departments {
		if err := db.Create(&department).Error; err != nil {
			return err
		}
	}
	return nil
}
