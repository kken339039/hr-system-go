package seeds

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type Seed struct {
	Name string
	Exec func(*gorm.DB) error
}

var Seeds []Seed

func CreateFile(name string) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s-%s.go", timestamp, name)
	filePath := filepath.Join("database", "seeds", fileName)

	// Create seeds directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create seeds directory: %w", err)
	}

	// Create the seed file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create seed file: %w", err)
	}
	defer file.Close()

	// Write the seed template to the file
	seedContent := fmt.Sprintf(`package seeds

import "gorm.io/gorm"

func init() {
		Seeds = append(Seeds, Seed{
        Name:      "%s-%s",
        Exec:       Exec_%s,
    })
}

func Exec_%s(db *gorm.DB) error {
    // TODO: Implement the seed logic here
    return nil
}
`, timestamp, name, timestamp, timestamp)

	_, err = file.WriteString(seedContent)
	if err != nil {
		return "", fmt.Errorf("failed to write seed content: %w", err)
	}

	return fileName, nil
}

func RunAll(db *gorm.DB) error {
	for _, seed := range Seeds {
		if err := runSeed(db, seed); err != nil {
			return err
		}
	}
	return nil
}

func Run(db *gorm.DB, name string) error {
	for _, seeder := range Seeds {
		if seeder.Name == name {
			return runSeed(db, seeder)
		}
	}
	return fmt.Errorf("seeder not found: %s", name)
}

func runSeed(db *gorm.DB, seed Seed) error {
	fmt.Printf("Running seeder: %s\n", seed.Name)
	if err := seed.Exec(db); err != nil {
		return fmt.Errorf("error running seeder %s: %w", seed.Name, err)
	}
	fmt.Printf("Seeder completed: %s\n", seed.Name)
	return nil
}
