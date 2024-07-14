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

type LeaveController struct {
	logger      *logger.Logger
	service     services.LeaveServiceInterface
	authService auth_service.AuthServiceInterface
}

func NewLeaveController(logger *logger.Logger, service services.LeaveServiceInterface, authService auth_service.AuthServiceInterface) *LeaveController {
	return &LeaveController{
		logger:      logger,
		service:     service,
		authService: authService,
	}
}

func (c *LeaveController) RegisterRoutes(r *gin.Engine) {
	leaveRoutes := r.Group("/api/users/:userId/leave")
	{
		leaveRoutes.GET("", c.authService.AuthUserAbilityWrapper(c.listLeaves, constants.ABILITY_READ_LEAVE))
		leaveRoutes.GET(":id", c.authService.AuthUserAbilityWrapper(c.getLeave, constants.ABILITY_READ_LEAVE))
		leaveRoutes.POST("", c.authService.AuthUserAbilityWrapper(c.createLeave, constants.ABILITY_READ_WRITE_LEAVE))
		leaveRoutes.PUT(":id", c.authService.AuthUserAbilityWrapper(c.updateLeave, constants.ABILITY_READ_WRITE_LEAVE))
		leaveRoutes.DELETE(":id", c.authService.AuthUserAbilityWrapper(c.deleteLeave, constants.ABILITY_DELETE_LEAVE))
	}
}

func (c *LeaveController) listLeaves(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Get Leaves"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	if !c.authService.AbleToAccessOtherUserData(ctx, userID, constants.ABILITY_ALL_GRANTS_LEAVE) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}

	pagination := utils.NewPagination(ctx)
	leaves, totalRows, err := c.service.FindLeavesByUserID(userID, &pagination)
	if err != nil {
		c.logger.Error("Failed to Find User", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewLeaveListResponse(leaves, totalRows, pagination))
}

func (c *LeaveController) getLeave(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Get Leave"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	currentUser := c.authService.GetCurrentUser(ctx)
	// only self read one leave
	if userID != int(currentUser.ID) {
		c.logger.Error("Cannot create leave which not belong currentUser")
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}

	leaveId := ctx.Param("id")
	leaveID, err := strconv.Atoi(leaveId)
	if err != nil {
		c.logger.Error("Cannot not parse Leave ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	leave, err := c.service.FindLeaveByID(leaveID)
	if err != nil {
		c.logger.Error("Cannot not find leave", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leave": dtos.NewLeaveResponse(leave)})
}

func (c *LeaveController) createLeave(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Create Leave"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	currentUser := c.authService.GetCurrentUser(ctx)
	// only self create leave
	if userID != int(currentUser.ID) {
		c.logger.Error("Cannot create leave which not belong currentUser")
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}
	var payload dtos.CreateLeaveRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot not parse create payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	leave, err := c.service.CreateLeaveByUser(currentUser, payload)
	if err != nil {
		c.logger.Error("Cannot not create leave", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leave": dtos.NewLeaveResponse(leave)})
}

func (c *LeaveController) updateLeave(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Update Leave"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	currentUser := c.authService.GetCurrentUser(ctx)
	// only self update leave
	if userID != int(currentUser.ID) {
		c.logger.Error("Cannot update leave which not belong currentUser")
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}
	var payload dtos.UpdateLeaveRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot not parse update payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	leaveId := ctx.Param("id")
	leaveID, err := strconv.Atoi(leaveId)
	if err != nil {
		c.logger.Error("Cannot not parse Leave ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	leave, err := c.service.UpdateLeaveByID(leaveID, payload)
	if err != nil {
		c.logger.Error("Cannot not update leave", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leave": dtos.NewLeaveResponse(leave)})
}

func (c *LeaveController) deleteLeave(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Delete Leave"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	currentUser := c.authService.GetCurrentUser(ctx)
	// only self delete leave
	if userID != int(currentUser.ID) {
		c.logger.Error("Cannot Delete leave which not belong currentUser")
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}
	leaveId := ctx.Param("id")
	leaveID, err := strconv.Atoi(leaveId)
	if err != nil {
		c.logger.Error("Cannot not parse Leave ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	if err := c.service.DeleteLeaveByID(leaveID); err != nil {
		c.logger.Error("Cannot not delete leave", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
