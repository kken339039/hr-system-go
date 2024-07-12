package migrations

import (
	"hr-system-go/internal/user/models"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create-user",
		Timestamp: "20240709211815",
		Up:        Up_20240709211815,
		Down:      Down_20240709211815,
	})
}

func Up_20240709211815(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

func Down_20240709211815(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.User{})
}
