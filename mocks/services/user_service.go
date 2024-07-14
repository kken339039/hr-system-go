package services

import (
	"hr-system-go/internal/user/dtos"
	"hr-system-go/internal/user/models"
	"hr-system-go/utils"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user *models.User, password string) error {
	args := m.Called(user, password)
	return args.Error(0)
}

func (m *MockUserService) FindUsers(pagination *utils.Pagination) ([]models.User, int64, error) {
	args := m.Called(pagination)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) FindUserByID(userId int) (*models.User, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUserByID(userId int, payload dtos.UpdateUserRequest) (*models.User, error) {
	args := m.Called(userId, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUserByID(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockUserService) UpdatePassword(user *models.User, newPassword string) error {
	args := m.Called(user, newPassword)
	return args.Error(0)
}

func (m *MockUserService) changeUserDepartment(user *models.User, newDeploymentId *int) error {
	args := m.Called(user, newDeploymentId)
	return args.Error(0)
}
