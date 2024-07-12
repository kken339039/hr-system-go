package migrations

import (
	"gorm.io/gorm"

	models "hr-system-go/internal/attendance/models"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create_leave",
		Timestamp: "20240710143647",
		Up:        Up_20240710143647,
		Down:      Down_20240710143647,
	})
}

func Up_20240710143647(db *gorm.DB) error {
	return db.AutoMigrate(&models.Leave{})
}

func Down_20240710143647(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Leave{})
}
