package session

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/session/controllers"

	"github.com/gin-gonic/gin"
)

type SessionModule struct {
	app.AppModuleInterface
}

func (m *SessionModule) Controllers() []interface{} {
	return []interface{}{
		controllers.NewSessionsController,
		func(
			r *gin.Engine,
			c *controllers.SessionsController,
			logger *logger.Logger,
		) *SessionModule {
			c.RegisterRoutes(r)
			logger.Info("= Session module init")
			return m
		},
	}
}

func (m *SessionModule) Provide() []interface{} {
	return []interface{}{}
}
