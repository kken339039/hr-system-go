package main

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/http"
	"hr-system-go/app/plugins/logger"
	user "hr-system-go/internal/users"
	// "net/http"
)

func main() {
	app := app.NewApplication()
	app.AddModule(&user.UserModule{})

	app.Run(func(
		env *env.Env,
		logger *logger.Logger,
		http *http.HttpServer,
		UserModule *user.UserModule,
	) {
		env.SetDefaultEnv(map[string]string{
			"PORT":        "3000",
			"ENVIRONMENT": "development",
		})
	})
}
