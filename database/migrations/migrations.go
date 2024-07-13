package migrations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration represents a single database migration record
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Timestamp string `gorm:"type:varchar(14);uniqueIndex"`
}

// MigrationFunc represents a function that performs a migration
type MigrationFunc func(*gorm.DB) error

// MigrationPair contains Up and Down functions for a migration
type MigrationPair struct {
	Name      string
	Timestamp string
	Up        MigrationFunc
	Down      MigrationFunc
}

// Migrations is a slice of all migrations
var Migrations []MigrationPair

// CreateMigrationFile creates a new migration file with the given name
func CreateFile(name string) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.go", timestamp, name)
	filePath := filepath.Join("database", "migrations", fileName)

	// Create migrations directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create the migration file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create migration file: %w", err)
	}
	defer file.Close()

	// Write the migration template to the file
	migrationContent := fmt.Sprintf(`package migrations

import "gorm.io/gorm"

func init() {
    Migrations = append(Migrations, MigrationPair{
        Name:      "%s",
        Timestamp: "%s",
        Up:        Up_%s,
        Down:      Down_%s,
    })
}

func Up_%s(db *gorm.DB) error {
    // TODO: Implement the migration logic here
    return nil
}

func Down_%s(db *gorm.DB) error {
    // TODO: Implement the rollback logic here
    return nil
}
`, name, timestamp, timestamp, timestamp, timestamp, timestamp)

	_, err = file.WriteString(migrationContent)
	if err != nil {
		return "", fmt.Errorf("failed to write migration content: %w", err)
	}

	return fileName, nil
}

// RunMigrations runs all pending migrations
func Run(db *gorm.DB) error {
	// Ensure the migration table exists
	err := db.AutoMigrate(&MigrationRecord{})
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get the last applied migration timestamp
	var lastMigration MigrationRecord
	if err := db.Order("timestamp desc").First(&lastMigration).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Sort migrations by timestamp
	sort.Slice(Migrations, func(i, j int) bool {
		return Migrations[i].Timestamp < Migrations[j].Timestamp
	})

	// Run migrations that are newer than the last applied migration
	for _, migration := range Migrations {
		if migration.Timestamp > lastMigration.Timestamp {
			fmt.Printf("Running migration: %s\n", migration.Name)
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
			}
			// Record the migration
			db.Create(&MigrationRecord{Timestamp: migration.Timestamp})
			fmt.Printf("Migration completed: %s\n", migration.Name)
		}
	}
	return nil
}

// RollbackMigration rolls back the last migration
func Rollback(db *gorm.DB) error {
	var lastMigration MigrationRecord
	if err := db.Order("timestamp desc").First(&lastMigration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	for i := len(Migrations) - 1; i >= 0; i-- {
		if Migrations[i].Timestamp == lastMigration.Timestamp {
			fmt.Printf("Rolling back migration: %s\n", Migrations[i].Name)
			if err := Migrations[i].Down(db); err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", Migrations[i].Name, err)
			}
			// Remove the migration record
			db.Delete(&lastMigration)
			fmt.Printf("Rollback completed: %s\n", Migrations[i].Name)
			return nil
		}
	}

	return fmt.Errorf("migration %s not found", lastMigration.Timestamp)
}
