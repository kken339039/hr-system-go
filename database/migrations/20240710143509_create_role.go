package migrations

import (
	"hr-system-go/internal/auth/models"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create_role",
		Timestamp: "20240710143509",
		Up:        Up_20240710143509,
		Down:      Down_20240710143509,
	})
}

func Up_20240710143509(db *gorm.DB) error {
	return db.AutoMigrate(&models.Role{})
}

func Down_20240710143509(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Role{})
}
