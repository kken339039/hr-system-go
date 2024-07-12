package main

import (
	"fmt"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/database/migrations"
	"hr-system-go/database/seeds"
	"os"
	"strings"

	"go.uber.org/zap"
)

func main() {
	env := env.NewEnv()
	logger := logger.NewLogger(env)

	if len(os.Args) < 2 {
		logger.Error(fmt.Sprintf("Usage: %s <command> [args...]", os.Args[0]))
		return
	}

	switch os.Args[1] {
	case "init":
		db := DBConnect(env, logger)
		defer db.Close()
		err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", env.GetEnv("DB_DATABASE")))
		if err != nil {
			logger.Error("Failed to Exec Command", zap.Error(err))
		}
		logger.Error("Database created or already exists")
	case "migration:create":
		if len(os.Args) != 3 {
			logger.Error(fmt.Sprintf("Usage: %s migration:create <migration-name>", os.Args[0]))
			return
		}
		migrateName, err := migrations.CreateFile(os.Args[2])
		if err != nil {
			logger.Error("Failed to create migration file:", zap.Error(err))
			return
		}
		logger.Info("Created migration file", zap.String("MigrateName", migrateName))
	case "migration:run":
		db := DBConnect(env, logger)
		defer db.Close()
		if err := migrations.Run(db.DB()); err != nil {
			logger.Error("Failed to run migrations", zap.Error(err))
			return
		}
		logger.Info("Run migrations successfully")
	case "migration:rollback":
		db := DBConnect(env, logger)
		defer db.Close()
		if err := migrations.Rollback(db.DB()); err != nil {
			logger.Error("Failed to rollback migration: %v", zap.Error(err))
			return
		}
		logger.Info("Rollback migration successfully")
	case "seed:create":
		if len(os.Args) != 3 {
			logger.Error(fmt.Sprintf("Usage: %s seed:create <seed-name>", os.Args[0]))
		}
		name, err := seeds.CreateFile(os.Args[2])
		if err != nil {
			logger.Error("Failed to create seed file:", zap.Error(err))
			return
		}
		logger.Info("Created migration file", zap.String("MigrateName", name))
	case "seed:runAll":
		db := DBConnect(env, logger)
		defer db.Close()
		if err := seeds.RunAll(db.DB()); err != nil {
			logger.Error("Failed to run all seeds", zap.Error(err))
			return
		}
		logger.Info("All seeders completed successfully")
	case "seed:run":
		db := DBConnect(env, logger)
		defer db.Close()
		if len(os.Args) != 3 {
			logger.Error(fmt.Sprintf("Usage: %s seed:run <seed-name>", os.Args[0]))
		}
		filename := strings.TrimSuffix(os.Args[2], ".go")
		if err := seeds.Run(db.DB(), filename); err != nil {
			logger.Error("Failed to run seed:", zap.Error(err))
			return
		}
		logger.Info(fmt.Sprintf("Run Seed %s completed successfully", filename))
	default:
		logger.Error(fmt.Sprintf("Unknown command: %s", os.Args[1]))
	}
}

func DBConnect(env *env.Env, logger *logger.Logger) *mysql.MySqlStore {
	store := mysql.NewMySqlStore(env, logger)
	store.Connect(
		env.GetEnv("DB_USER"),
		env.GetEnv("DB_PASSWORD"),
		env.GetEnv("DB_DATABASE"),
		env.GetEnv("DB_HOST"),
		env.GetEnv("DB_PORT"),
		env.GetEnv("DB_PARAMS"),
	)
	return store
}
