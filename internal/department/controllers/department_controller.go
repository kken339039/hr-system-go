package controllers

import (
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/auth/constants"
	auth_service "hr-system-go/internal/auth/services"
	"hr-system-go/internal/department/dtos"
	"hr-system-go/internal/department/services"
	"hr-system-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DepartmentController struct {
	logger      *logger.Logger
	service     services.DepartmentServiceInterface
	authService auth_service.AuthServiceInterface
}

func NewDepartmentController(logger *logger.Logger, service services.DepartmentServiceInterface, authService auth_service.AuthServiceInterface) *DepartmentController {
	return &DepartmentController{
		logger:      logger,
		service:     service,
		authService: authService,
	}
}

func (c *DepartmentController) RegisterRoutes(r *gin.Engine) {
	departmentRoutes := r.Group("/api/department")
	{
		departmentRoutes.GET("", c.authService.AuthUserAbilityWrapper(c.listDepartments, constants.ABILITY_READ_DEPARTMENT))
		departmentRoutes.GET("/:id", c.authService.AuthUserAbilityWrapper(c.GetDepartment, constants.ABILITY_READ_DEPARTMENT))
		departmentRoutes.POST("", c.authService.AuthUserAbilityWrapper(c.CreateDepartment, constants.ABILITY_READ_DEPARTMENT))
		departmentRoutes.PUT("/:id", c.authService.AuthUserAbilityWrapper(c.UpdateDepartment, constants.ABILITY_READ_WRITE_DEPARTMENT))
		departmentRoutes.DELETE("/:id", c.authService.AuthUserAbilityWrapper(c.DeleteDepartment, constants.ABILITY_DELETE_DEPARTMENT))
	}
}

func (c *DepartmentController) listDepartments(ctx *gin.Context) {
	pagination := utils.NewPagination(ctx)
	departments, totalRows, err := c.service.FindDepartments(&pagination)
	if err != nil {
		c.logger.Error("Failed to Find User", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Find Departments Error"})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewDepartmentListResponse(departments, totalRows, pagination))
}

func (c *DepartmentController) GetDepartment(ctx *gin.Context) {
	departmentId := ctx.Param("id")
	departmentID, err := strconv.Atoi(departmentId)
	errorMsg := "Failed to Get Department"
	if err != nil {
		c.logger.Error("Cannot not parse Department ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	department, err := c.service.FindDepartmentByID(departmentID)
	if err != nil {
		c.logger.Error("Failed to Create User", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewDepartmentResponse(department))
}

func (c *DepartmentController) CreateDepartment(ctx *gin.Context) {
	var payload dtos.CreateDepartmentRequest
	errorMsg := "Failed to Create Department"

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot not parse create payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	department, err := c.service.CreateDepartment(payload)
	if err != nil {
		c.logger.Error("Cannot not create user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewDepartmentResponse(department))
}

func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
	departmentId := ctx.Param("id")
	departmentID, err := strconv.Atoi(departmentId)
	errorMsg := "Failed to Update Department"
	if err != nil {
		c.logger.Error("Cannot not parse Department ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	var payload dtos.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot not parse update payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	department, err := c.service.UpdateDepartmentByID(departmentID, payload)
	if err != nil {
		c.logger.Error("Cannot not update user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusOK, dtos.NewDepartmentResponse(department))
}

func (c *DepartmentController) DeleteDepartment(ctx *gin.Context) {
	departmentId := ctx.Param("id")
	departmentID, err := strconv.Atoi(departmentId)
	errorMsg := "Failed to Delete Department"
	if err != nil {
		c.logger.Error("Cannot not parse Department ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	if err := c.service.DeleteDepartmentByID(departmentID); err != nil {
		c.logger.Error("Cannot not delete Department", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
