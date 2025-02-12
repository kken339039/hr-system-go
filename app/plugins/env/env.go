package env

import (
	"fmt"
	"log"
	"os"

	"hr-system-go/app/plugins"

	"github.com/joho/godotenv"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewEnv)
}

type Env struct {
	defaultValues map[string]string
}

func (e *Env) SetDefaultEnv(values map[string]string) {
	for key, value := range values {
		e.defaultValues[key] = value
	}
}

func (e Env) GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return e.defaultValues[key]
}

func NewEnv() *Env {
	env := &Env{
		defaultValues: map[string]string{
			"ENVIRONMENT": "development",
		},
	}
	// for test with different package
	projectRoot := os.Getenv("PROJECT_ROOT")
	if len(projectRoot) == 0 {
		projectRoot, _ = os.Getwd()
	}

	path := fmt.Sprintf("%s/.env", projectRoot)

	if os.Getenv("ENVIRONMENT") == "test" {
		path = fmt.Sprintf("%s.test", path)
	}

	if err := godotenv.Load(path); err != nil {
		log.Print("No .env file found", err)
	}

	return env
}
