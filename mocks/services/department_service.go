package services

import (
	"hr-system-go/internal/department/dtos"
	"hr-system-go/internal/department/models"
	"hr-system-go/utils"

	"github.com/stretchr/testify/mock"
)

type MockDepartmentService struct {
	mock.Mock
}

func (m *MockDepartmentService) FindDepartments(pagination *utils.Pagination) ([]models.Department, int64, error) {
	args := m.Called(pagination)
	return args.Get(0).([]models.Department), args.Get(1).(int64), args.Error(2)
}

func (m *MockDepartmentService) FindDepartmentByID(id int) (*models.Department, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentService) CreateDepartment(payload dtos.CreateDepartmentRequest) (*models.Department, error) {
	args := m.Called(payload)
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentService) UpdateDepartmentByID(id int, payload dtos.UpdateDepartmentRequest) (*models.Department, error) {
	args := m.Called(id, payload)
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentService) DeleteDepartmentByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
