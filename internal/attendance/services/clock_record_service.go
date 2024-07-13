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

	err := s.db.Model(&models.ClockRecord{}).Preload("User").Limit(pagination.Limit).Offset(pagination.Offset()).Order(pagination.Sort).Find(&records, "user_id = ?", userID).Error
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
