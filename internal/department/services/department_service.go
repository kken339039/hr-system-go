package services

import (
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"

	"hr-system-go/internal/department/dtos"
	"hr-system-go/internal/department/models"
	"hr-system-go/utils"

	"go.uber.org/zap"
)

type DepartmentService struct {
	logger *logger.Logger
	db     *mysql.MySqlStore
}

func NewDepartmentService(logger *logger.Logger, db *mysql.MySqlStore) *DepartmentService {
	return &DepartmentService{
		logger: logger,
		db:     db,
	}
}

func (s *DepartmentService) FindDepartments(pagination utils.Pagination) ([]models.Department, int64, error) {
	var departments []models.Department
	var totalCount int64 = 0

	if err := models.ValidScope(s.db.DB()).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := models.ValidScope(s.db.DB()).Limit(pagination.Limit).Offset(pagination.Offset()).Order(pagination.Sort).Find(&departments).Error
	if err != nil {
		return nil, 0, err
	}

	return departments, totalCount, nil
}

func (s *DepartmentService) FindDepartmentByID(departmentID int) (*models.Department, error) {
	var department *models.Department
	if err := models.ValidScope(s.db.DB()).First(&department, departmentID).Error; err != nil {
		s.logger.Error("Cannot Not Find User by ID", zap.Error(err))
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) CreateDepartmentByID(payload dtos.CreateDepartmentRequest) (*models.Department, error) {
	department := &models.Department{Name: payload.Name}
	if payload.Descriptions != nil {
		department.Descriptions = *payload.Descriptions
	}
	if err := models.ValidScope(s.db.DB()).Create(&department).Error; err != nil {
		s.logger.Error("Cannot Update Deployment Data", zap.Error(err))
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) UpdateDepartmentByID(departmentID int, payload dtos.UpdateDepartmentRequest) (*models.Department, error) {
	var department *models.Department
	if err := models.ValidScope(s.db.DB()).First(&department, departmentID).Updates(payload).Error; err != nil {
		s.logger.Error("Cannot Update Deployment Data", zap.Error(err))
		return nil, err
	}

	return s.FindDepartmentByID(departmentID)
}

func (s *DepartmentService) DeleteDepartmentByID(departmentID int) error {
	var department *models.Department
	return models.ValidScope(s.db.DB()).First(&department, departmentID).Update("status", "removed").Error
}
