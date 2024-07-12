package seeds

import (
	user_models "hr-system-go/internal/user/models"

	"gorm.io/gorm"
)

func init() {
	Seeds = append(Seeds, Seed{
		Name: "20240710120332-import-users",
		Exec: Exec_20240710120332,
	})
}

func Exec_20240710120332(db *gorm.DB) error {
	users := []user_models.User{
		{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			Age:   22,
		},
		{
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
			Age:   30,
		},
		{
			Name:  "Bob Johnson",
			Email: "bob.johnson@example.com",
			Age:   35,
		},
		{
			Name:  "Alice Brown",
			Email: "alice.brown@example.com",
			Age:   24,
		},
		{
			Name:  "Charlie Wilson",
			Email: "charlie.wilson@example.com",
			Age:   35,
		},
	}
	for _, user := range users {
		user.GenerateRandomPassword()
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}
	return nil
}
