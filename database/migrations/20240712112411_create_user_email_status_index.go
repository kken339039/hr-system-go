package migrations

import (
	user_models "hr-system-go/internal/user/models"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, MigrationPair{
		Name:      "create_user_email_status_index",
		Timestamp: "20240712112411",
		Up:        Up_20240712112411,
		Down:      Down_20240712112411,
	})
}

func Up_20240712112411(db *gorm.DB) error {
	migrator := db.Migrator()
	if !migrator.HasIndex(&user_models.User{}, "idx_email_status") {
		err := migrator.CreateIndex(&user_models.User{}, "idx_email_status")
		if err != nil {
			return err
		}
	}
	return nil
}

func Down_20240712112411(db *gorm.DB) error {
	err := db.Migrator().DropIndex(&user_models.User{}, "idx_email_status")
	return err
}
