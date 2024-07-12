package attendance

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/attendance/controllers"
	"hr-system-go/internal/attendance/services"

	"github.com/gin-gonic/gin"
)

type AttendanceModule struct {
	app.AppModuleInterface
}

func (m *AttendanceModule) Controllers() []interface{} {
	return []interface{}{
		controllers.NewLeaveController,
		func(
			r *gin.Engine,
			lc *controllers.LeaveController,
			logger *logger.Logger,
		) *AttendanceModule {
			lc.RegisterRoutes(r)
			logger.Info("= Attendance module init")
			return m
		},
	}
}

func (m *AttendanceModule) Provide() []interface{} {
	return []interface{}{
		services.NewLeaveService,
	}
}
