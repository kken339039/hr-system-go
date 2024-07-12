package migrations

import (
	"hr-system-go/internal/department/models"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create_department",
		Timestamp: "20240710143709",
		Up:        Up_20240710143709,
		Down:      Down_20240710143709,
	})
}

func Up_20240710143709(db *gorm.DB) error {
	// TODO: Implement the migration logic here
	return db.AutoMigrate(&models.Department{})
}

func Down_20240710143709(db *gorm.DB) error {
	// TODO: Implement the rollback logic here
	return db.Migrator().DropTable(&models.Department{})
}
