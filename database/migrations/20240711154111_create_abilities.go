package migrations

import (
	"hr-system-go/internal/auth/models"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create_abilities",
		Timestamp: "20240711154111",
		Up:        Up_20240711154111,
		Down:      Down_20240711154111,
	})
}

func Up_20240711154111(db *gorm.DB) error {
	// TODO: Implement the migration logic here
	return db.AutoMigrate(&models.Ability{})
}

func Down_20240711154111(db *gorm.DB) error {
	// TODO: Implement the rollback logic here
	return db.Migrator().DropTable(&models.Ability{})
}
