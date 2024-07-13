package http

import (
	"hr-system-go/app/plugins"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/http/interceptors"
	"hr-system-go/app/plugins/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewRouter)
}

func NewRouter(env *env.Env, logger *logger.Logger) *gin.Engine {
	mode := env.GetEnv("ENVIRONMENT")
	var ginMode string
	switch strings.ToLower(mode) {
	case "development":
		ginMode = gin.DebugMode
	case "staging":
		ginMode = gin.DebugMode
	case "preview":
		ginMode = gin.ReleaseMode
	case "production":
		ginMode = gin.ReleaseMode
	case "test":
		ginMode = gin.TestMode
	default:
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	r := gin.New()
	r.Use(gin.Recovery(), interceptors.RequestLog(logger))

	return r
}
