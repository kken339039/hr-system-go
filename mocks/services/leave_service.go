package services

import (
	"hr-system-go/internal/attendance/dtos"
	"hr-system-go/internal/attendance/models"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/utils"

	"github.com/stretchr/testify/mock"
)

type MockLeaveService struct {
	mock.Mock
}

func (m *MockLeaveService) FindLeavesByUserID(userID int, pagination *utils.Pagination) ([]models.Leave, int64, error) {
	args := m.Called(userID, pagination)
	return args.Get(0).([]models.Leave), args.Get(1).(int64), args.Error(2)
}

func (m *MockLeaveService) FindLeaveByID(leaveID int) (*models.Leave, error) {
	args := m.Called(leaveID)
	return args.Get(0).(*models.Leave), args.Error(1)
}

func (m *MockLeaveService) CreateLeaveByUser(user *user_models.User, payload dtos.CreateLeaveRequest) (*models.Leave, error) {
	args := m.Called(user, payload)
	return args.Get(0).(*models.Leave), args.Error(1)
}

func (m *MockLeaveService) UpdateLeaveByID(leaveID int, payload dtos.UpdateLeaveRequest) (*models.Leave, error) {
	args := m.Called(leaveID, payload)
	return args.Get(0).(*models.Leave), args.Error(1)
}

func (m *MockLeaveService) DeleteLeaveByID(leaveID int) error {
	args := m.Called(leaveID)
	return args.Error(0)
}
