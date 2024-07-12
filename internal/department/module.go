package department

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/department/controllers"
	"hr-system-go/internal/department/services"

	"github.com/gin-gonic/gin"
)

type DepartmentModule struct {
	app.AppModuleInterface
}

func (m *DepartmentModule) Controllers() []interface{} {
	return []interface{}{
		controllers.NewDepartmentController,
		func(
			r *gin.Engine,
			c *controllers.DepartmentController,
			logger *logger.Logger,
		) *DepartmentModule {
			c.RegisterRoutes(r)
			logger.Info("= Department module init")
			return m
		},
	}
}

func (m *DepartmentModule) Provide() []interface{} {
	return []interface{}{
		services.NewDepartmentService,
	}
}
