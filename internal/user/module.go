package user

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/user/controllers"
	"hr-system-go/internal/user/services"

	"github.com/gin-gonic/gin"
)

type UserModule struct {
	app.AppModuleInterface
}

func (m *UserModule) Controllers() []interface{} {
	return []interface{}{
		controllers.NewUsersController,
		func(
			r *gin.Engine,
			c *controllers.UsersController,
			logger *logger.Logger,
		) *UserModule {
			c.RegisterRoutes(r)
			logger.Info("= User module init")
			return m
		},
	}
}

func (m *UserModule) Provide() []interface{} {
	return []interface{}{
		services.NewUserService,
	}
}
