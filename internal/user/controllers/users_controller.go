package controllers

import (
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/auth/constants"
	auth_service "hr-system-go/internal/auth/services"
	dtos "hr-system-go/internal/user/dtos"
	"hr-system-go/internal/user/services"
	"hr-system-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UsersController struct {
	logger      *logger.Logger
	service     services.UserServiceInterface
	authService auth_service.AuthServiceInterface
}

func NewUsersController(logger *logger.Logger, service services.UserServiceInterface, authService auth_service.AuthServiceInterface) *UsersController {
	return &UsersController{
		logger:      logger,
		service:     service,
		authService: authService,
	}
}

func (c *UsersController) RegisterRoutes(r *gin.Engine) {
	userRoutes := r.Group("/api/users")
	{
		userRoutes.GET("", c.authService.AuthUserAbilityWrapper(c.listUsers, constants.ABILITY_ALL_GRANTS_USER))
		userRoutes.GET("/:userId", c.authService.AuthUserAbilityWrapper(c.GetUser, constants.ABILITY_READ_USER))
		userRoutes.PUT("/:userId", c.authService.AuthUserAbilityWrapper(c.UpdateUser, constants.ABILITY_READ_WRITE_USER))
		userRoutes.DELETE("/:userId", c.authService.AuthUserAbilityWrapper(c.DeleteUser, constants.ABILITY_DELETE_USER))
	}
}

func (c *UsersController) listUsers(ctx *gin.Context) {
	pagination := utils.NewPagination(ctx)
	users, totalRows, err := c.service.FindUsers(&pagination)
	if err != nil {
		c.logger.Error("Failed to Find User", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Find users Error"})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewUserListResponse(users, totalRows, pagination))
}

func (c *UsersController) GetUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Get User"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	if !c.authService.AbleToAccessOtherUserData(ctx, userID, constants.ABILITY_ALL_GRANTS_USER) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}

	user, err := c.service.FindUserByID(userID)
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewUserResponse(user))
}

func (c *UsersController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Update User"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	if !c.authService.AbleToAccessOtherUserData(ctx, userID, constants.ABILITY_ALL_GRANTS_USER) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}

	var payload dtos.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot not parse update payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	user, err := c.service.UpdateUserByID(userID, payload)
	if err != nil {
		c.logger.Error("Cannot not update user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewUserResponse(user))
}

func (c *UsersController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userID, err := strconv.Atoi(userId)
	errorMsg := "Failed to Delete User"
	if err != nil {
		c.logger.Error("Cannot not parse User ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	if !c.authService.AbleToAccessOtherUserData(ctx, userID, constants.ABILITY_ALL_GRANTS_USER) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})
		return
	}

	if err := c.service.DeleteUserByID(userID); err != nil {
		c.logger.Error("Cannot not delete user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
