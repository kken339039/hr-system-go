package migrations

import (
	"hr-system-go/internal/attendance/models"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create_clock_record",
		Timestamp: "20240710143654",
		Up:        Up_20240710143654,
		Down:      Down_20240710143654,
	})
}

func Up_20240710143654(db *gorm.DB) error {
	// TODO: Implement the migration logic here
	return db.AutoMigrate(&models.ClockRecord{})
}

func Down_20240710143654(db *gorm.DB) error {
	// TODO: Implement the rollback logic here
	return db.Migrator().DropTable(&models.ClockRecord{})
}
