package services

import (
	"errors"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/internal/attendance/models"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ClockRecordService struct {
	logger *logger.Logger
	db     *mysql.MySqlStore
}

func NewClockRecordService(logger *logger.Logger, db *mysql.MySqlStore) *ClockRecordService {
	return &ClockRecordService{
		logger: logger,
		db:     db,
	}
}

func (s ClockRecordService) FindClockRecordsByUserID(userID int, pagination utils.Pagination) ([]models.ClockRecord, int64, error) {
	var records []models.ClockRecord
	var totalCount int64 = 0

	if err := s.db.Model(&models.ClockRecord{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.Model(&models.ClockRecord{}).Preload("User").Limit(pagination.Limit).Offset(pagination.Offset()).Order(pagination.Sort).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, totalCount, nil
}

func (s *ClockRecordService) ClockByUser(user *user_models.User) (*models.ClockRecord, error) {
	var existRecord *models.ClockRecord
	recordBaseQuery := s.db.DB().Preload("User").Where(&models.ClockRecord{UserID: user.ID}).Where("clock_out is NULL")
	result := recordBaseQuery.First(&existRecord)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newRecord := &models.ClockRecord{
			User:    *user,
			ClockIn: time.Now(),
		}
		if err := s.db.DB().Create(&newRecord).Error; err != nil {
			s.logger.Error("Clock In Failed", zap.Error(err))
			return nil, err
		}
		return newRecord, nil
	} else {
		if err := recordBaseQuery.First(&existRecord).Update("clock_out", time.Now()).Error; err != nil {
			s.logger.Error("Clock Out Failed", zap.Error(err))
			return nil, err
		}
		return existRecord, nil
	}
}

// func (s *ClockRecordService) UpdateLeaveByID(leaveID int, payload dtos.UpdateLeaveRequest) (*models.Leave, error) {
// 	leave := &models.Leave{}

// 	if payload.LeaveType != nil {
// 		leave.LeaveType = *payload.LeaveType
// 	}
// 	if payload.Status != nil {
// 		leave.Status = *payload.Status
// 	}

// 	if payload.StartDate != nil {
// 		startDate, _ := utils.ParseDateTime(*payload.StartDate)
// 		leave.StartDate = startDate
// 	}

// 	if payload.EndDate != nil {
// 		endDate, _ := utils.ParseDateTime(*payload.EndDate)
// 		leave.EndDate = endDate
// 	}

// 	var updatedLeave *models.Leave
// 	if err := models.ValidLeaveScope(s.db.DB()).First(&updatedLeave, leaveID).Updates(leave).Error; err != nil {
// 		s.logger.Error("Cannot Update Leave Data", zap.Error(err))
// 		return nil, err
// 	}

// 	return s.FindLeaveByID(leaveID)
// }

// func (s *ClockRecordService) DeleteLeaveByID(leaveID int) error {
// 	var leave *models.Leave
// 	if err := models.ValidLeaveScope(s.db.DB()).First(&leave, leaveID).Update("status", "removed").Error; err != nil {
// 		s.logger.Error("Cannot Delete User", zap.Error(err))
// 		return err
// 	}

// 	return nil
// }

// func (s *ClockRecordService) FindValidRecord(userID int) (*models.Leave, error) {
// 	var existLeave *models.Leave
// 	if err := s.db.Model(&models.Leave{}).Preload("User").First(&existLeave, leaveID).Error; err != nil {
// 		s.logger.Error("Cannot Not Find Leave by ID", zap.Error(err))
// 		return nil, err
// 	}

// 	return leave, nil
// }
