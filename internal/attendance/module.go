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
		controllers.NewClockRecordController,
		func(
			r *gin.Engine,
			lc *controllers.LeaveController,
			crc *controllers.ClockRecordController,
			logger *logger.Logger,
		) *AttendanceModule {
			lc.RegisterRoutes(r)
			crc.RegisterRoutes(r)
			logger.Info("= Attendance module init")
			return m
		},
	}
}

func (m *AttendanceModule) Provide() []interface{} {
	return []interface{}{
		services.NewLeaveService,
		services.NewClockRecordService,
	}
}
