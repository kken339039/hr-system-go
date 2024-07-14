package services

import (
	"hr-system-go/internal/attendance/models"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/utils"

	"github.com/stretchr/testify/mock"
)

type MockClockRecordService struct {
	mock.Mock
}

func (m *MockClockRecordService) FindClockRecordsByUserID(userID int, pagination *utils.Pagination) ([]models.ClockRecord, int64, error) {
	args := m.Called(userID, pagination)
	return args.Get(0).([]models.ClockRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockClockRecordService) ClockByUser(user *user_models.User) (*models.ClockRecord, error) {
	args := m.Called(user)
	return args.Get(0).(*models.ClockRecord), args.Error(1)
}
