package services

import (
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/internal/attendance/dtos"
	"hr-system-go/internal/attendance/models"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/utils"

	"go.uber.org/zap"
)

type LeaveServiceInterface interface {
	FindLeavesByUserID(userID int, pagination *utils.Pagination) ([]models.Leave, int64, error)
	FindLeaveByID(leaveID int) (*models.Leave, error)
	CreateLeaveByUser(user *user_models.User, payload dtos.CreateLeaveRequest) (*models.Leave, error)
	UpdateLeaveByID(leaveID int, payload dtos.UpdateLeaveRequest) (*models.Leave, error)
	DeleteLeaveByID(leaveID int) error
}

type LeaveService struct {
	logger *logger.Logger
	db     *mysql.MySqlStore
}

func NewLeaveService(logger *logger.Logger, db *mysql.MySqlStore) *LeaveService {
	return &LeaveService{
		logger: logger,
		db:     db,
	}
}

func (s LeaveService) FindLeavesByUserID(userID int, pagination utils.Pagination) ([]models.Leave, int64, error) {
	var leaves []models.Leave
	var totalCount int64 = 0

	if err := models.ValidLeaveScope(s.db.DB()).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := models.ValidLeaveScope(s.db.DB()).Preload("User").Limit(pagination.Limit).Offset(pagination.Offset()).Order(pagination.Sort).Find(&leaves, "user_id = ?", userID).Error
	if err != nil {
		return nil, 0, err
	}
	return leaves, totalCount, nil
}

func (s *LeaveService) CreateLeaveByUser(user *user_models.User, payload dtos.CreateLeaveRequest) (*models.Leave, error) {
	startDate, _ := utils.ParseDateTime(*payload.StartDate)
	endDate, _ := utils.ParseDateTime(*payload.EndDate)
	leave := &models.Leave{
		UserID:    user.ID,
		StartDate: startDate,
		EndDate:   endDate,
		LeaveType: *payload.LeaveType,
	}

	if err := s.db.DB().Create(&leave).Error; err != nil {
		s.logger.Error("Create Leave Failed", zap.Error(err))
		return nil, err
	}

	return s.FindLeaveByID(int(leave.ID))
}

func (s *LeaveService) UpdateLeaveByID(leaveID int, payload dtos.UpdateLeaveRequest) (*models.Leave, error) {
	leave := &models.Leave{}

	if payload.LeaveType != nil {
		leave.LeaveType = *payload.LeaveType
	}
	if payload.Status != nil {
		leave.Status = *payload.Status
	}

	if payload.StartDate != nil {
		startDate, _ := utils.ParseDateTime(*payload.StartDate)
		leave.StartDate = startDate
	}

	if payload.EndDate != nil {
		endDate, _ := utils.ParseDateTime(*payload.EndDate)
		leave.EndDate = endDate
	}

	var updatedLeave *models.Leave
	if err := models.ValidLeaveScope(s.db.DB()).First(&updatedLeave, leaveID).Updates(leave).Error; err != nil {
		s.logger.Error("Cannot Update Leave Data", zap.Error(err))
		return nil, err
	}

	return s.FindLeaveByID(leaveID)
}

func (s *LeaveService) DeleteLeaveByID(leaveID int) error {
	var leave *models.Leave
	if err := models.ValidLeaveScope(s.db.DB()).First(&leave, leaveID).Update("status", "removed").Error; err != nil {
		s.logger.Error("Cannot Delete User", zap.Error(err))
		return err
	}

	return nil
}

func (s *LeaveService) FindLeaveByID(leaveID int) (*models.Leave, error) {
	var leave *models.Leave
	if err := models.ValidLeaveScope(s.db.DB()).Preload("User").First(&leave, leaveID).Error; err != nil {
		s.logger.Error("Cannot Not Find Leave by ID", zap.Error(err))
		return nil, err
	}

	return leave, nil
}
