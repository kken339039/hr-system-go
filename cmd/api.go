package main

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/http"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/app/plugins/redis"
	"hr-system-go/internal/attendance"
	"hr-system-go/internal/auth"
	"hr-system-go/internal/department"
	"hr-system-go/internal/session"
	"hr-system-go/internal/user"

	"strconv"
)

func main() {
	app := app.NewApplication()
	app.AddModule(&auth.AuthModule{})
	app.AddModule(&user.UserModule{})
	app.AddModule(&attendance.AttendanceModule{})
	app.AddModule(&session.SessionModule{})
	app.AddModule(&department.DepartmentModule{})

	app.Run(func(
		env *env.Env,
		logger *logger.Logger,
		http *http.HttpServer,
		authModule *auth.AuthModule,
		userModule *user.UserModule,
		attendanceModule *attendance.AttendanceModule,
		sessionModule *session.SessionModule,
		departmentModule *department.DepartmentModule,
		mysql *mysql.MySqlStore,
		redis *redis.RedisStore,
	) {
		env.SetDefaultEnv(map[string]string{
			"PORT":        "3000",
			"ENVIRONMENT": "development",
		})
		// DB connect
		mysql.Connect(
			env.GetEnv("DB_USER"),
			env.GetEnv("DB_PASSWORD"),
			env.GetEnv("DB_DATABASE"),
			env.GetEnv("DB_HOST"),
			env.GetEnv("DB_PORT"),
			env.GetEnv("DB_PARAMS"),
		)
		// REDIS connect
		redisDB, _ := strconv.Atoi(env.GetEnv("REDIS_DB"))
		redis.Connect(
			env.GetEnv("REDIS_HOST"),
			env.GetEnv("REDIS_PORT"),
			redisDB,
		)
	})
}
