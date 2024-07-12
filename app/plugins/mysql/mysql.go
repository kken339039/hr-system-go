package mysql

import (
	"fmt"
	"hr-system-go/app/plugins"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewMySqlStore)
}

var Store *MySqlStore

type MySqlStore struct {
	env    *env.Env
	logger *logger.Logger
	db     *gorm.DB
}

type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	AppliedAt time.Time
}

func NewMySqlStore(env *env.Env, logger *logger.Logger) *MySqlStore {
	Store = &MySqlStore{
		env:    env,
		logger: logger,
	}
	return Store
}

func (s *MySqlStore) Connect(username string, password string, database string, host string, port string, params string) {
	var paramsString string
	if params != "" {
		paramsString = fmt.Sprintf("?%s", params)
	} else {
		paramsString = ""
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", username, password, host, port, database, paramsString)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		s.logger.Error("Failed to connect DB", zap.Error(err))
		panic(err)
	}

	s.db = db
	s.logger.Info("Connected to MySQL database")
}

func (s *MySqlStore) DB() *gorm.DB {
	return s.db
}

func (s *MySqlStore) Model(instance interface{}) *gorm.DB {
	return s.db.Model(instance)
}

func (s *MySqlStore) Exec(cmd string) error {
	result := s.db.Exec(cmd)
	return result.Error
}

func (s *MySqlStore) Close() {
	instance, _ := s.db.DB()
	instance.Close()
	s.logger.Info("Close database connection")
}
