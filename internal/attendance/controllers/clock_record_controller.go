package controllers

import (
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/attendance/dtos"
	"hr-system-go/internal/attendance/services"
	"hr-system-go/internal/auth/constants"
	auth_service "hr-system-go/internal/auth/services"
	"hr-system-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ClockRecordController struct {
	logger      *logger.Logger
	service     *services.ClockRecordService
	authService *auth_service.AuthService
}

func NewClockRecordController(logger *logger.Logger, service *services.ClockRecordService, authService *auth_service.AuthService) *ClockRecordController {
	return &ClockRecordController{
		logger:      logger,
		service:     service,
		authService: authService,
	}
}

func (c *ClockRecordController) RegisterRoutes(r *gin.Engine) {
	clockRecordRoutes := r.Group("/api/users/:userId/clockRecord")
	{
		clockRecordRoutes.GET("", c.authService.AuthUserAbilityWrapper(c.listClockRecord, constants.ABILITY_READ_CLOCK_RECORD))
		clockRecordRoutes.POST("/clock", c.authService.AuthUserAbilityWrapper(c.touchClockRecord, constants.ABILITY_READ_WRITE_LEAVE))
	}
}

func (c *ClockRecordController) listClockRecord(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Get Clock Records"})
		return
	}

	if !c.authService.VerifyAllGrants(ctx, userID, constants.ABILITY_ALL_GRANTS_CLOCK_RECORD) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to Get Clock Records"})
		return
	}

	pagination := utils.NewPagination(ctx)
	records, totalRows, err := c.service.FindClockRecordsByUserID(userID, pagination)
	if err != nil {
		c.logger.Error("Failed to Find User", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Find Clock Records Error"})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewClockRecordListResponse(records, totalRows, pagination))
}

func (c *ClockRecordController) touchClockRecord(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Touch Record"})
		return
	}
	currentUser := c.authService.GetCurrentUser(ctx)
	// only self can clock in/out
	if userID != int(currentUser.ID) {
		c.logger.Error("Cannot create Clock Record which not belong currentUser")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to Touch Record"})
		return
	}

	record, err := c.service.ClockByUser(currentUser)
	if err != nil {
		c.logger.Error("Cannot not touch ClockRecord", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Touch ClockRecord"})
		return
	}
	// ctx.JSON(http.StatusOK, nil)
	ctx.JSON(http.StatusOK, gin.H{"leave": dtos.NewClockRecordResponse(record)})
}
