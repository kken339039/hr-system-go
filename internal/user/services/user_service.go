package services

import (
	"errors"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	department_models "hr-system-go/internal/department/models"
	"hr-system-go/internal/user/dtos"
	"hr-system-go/internal/user/models"
	"hr-system-go/utils"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger *logger.Logger
	db     *mysql.MySqlStore
}

func NewUserService(logger *logger.Logger, db *mysql.MySqlStore) *UserService {
	return &UserService{
		logger: logger,
		db:     db,
	}
}

func (s *UserService) RegisterUser(user *models.User, password string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.PasswordEncrypt = string(hashedPassword)
	user.JoinDate = time.Now()

	if err := s.db.DB().Create(&user).Error; err != nil {
		s.logger.Error("Create User Failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserService) FindUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64 = 0

	if err := models.ValidScope(s.db.DB()).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := models.ValidScope(s.db.DB()).Limit(pagination.Limit).Offset(pagination.Offset()).Order(pagination.Sort).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, totalCount, nil
}

func (s *UserService) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := models.ValidScope(s.db.DB()).Where("email = ?", email).First(&user).Error; err != nil {
		s.logger.Error("Cannot Not Find User by email", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindUserByID(userId int) (*models.User, error) {
	var user *models.User
	if err := models.ValidScope(s.db.DB()).Preload("Role").Preload("Department").First(&user, userId).Error; err != nil {
		s.logger.Error("Cannot Not Find User by ID", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUserByID(userId int, payload dtos.UpdateUserRequest) (*models.User, error) {
	var user *models.User
	// var oldDepartmentID *uint
	if err := models.ValidScope(s.db.DB()).First(&user, userId).Updates(payload).Error; err != nil {
		s.logger.Error("Cannot Update User Data", zap.Error(err))
		return nil, err
	}

	if err := s.changeUserDepartment(user, payload.DepartmentID); err != nil {
		return nil, err
	}

	return s.FindUserByID(userId)
}

func (s *UserService) DeleteUserByID(userId int) error {
	var user *models.User
	if err := models.ValidScope(s.db.DB()).First(&user, userId).Update("status", "removed").Error; err != nil {
		s.logger.Error("Cannot Delete User", zap.Error(err))
		return err
	}

	if user.DepartmentID != nil {
		var department *department_models.Department
		if err := department_models.ValidScope(s.db.DB()).First(&department, &user.DepartmentID).Error; err != nil {
			s.logger.Error("Cannot Find User's Department", zap.Error(err))
			return err
		}
		if err := department.UpdateEmployCount(s.db.DB(), -1); err != nil {
			s.logger.Error("Cannot Update Old Department Employ Count", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *UserService) UpdatePassword(user *models.User, newPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordEncrypt), []byte(newPassword))
	// err nil means new password is same as old
	if err == nil {
		return errors.New("password is not changed")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = models.ValidScope(s.db.DB()).First(&user, user.ID).Update("PasswordEncrypt", hashedPassword).Error
	if err != nil {
		s.logger.Error("Cannot Update Password", zap.Error(err))
		return err
	}
	return nil
}

func (s *UserService) changeUserDepartment(user *models.User, newDeploymentId *int) error {
	if newDeploymentId != nil {
		var newDepartment *department_models.Department
		if err := department_models.ValidScope(s.db.DB()).First(&newDepartment, &newDeploymentId).Error; err != nil {
			s.logger.Error("Cannot Find Updating Department", zap.Error(err))
			return err
		}
		if err := newDepartment.UpdateEmployCount(s.db.DB(), 1); err != nil {
			s.logger.Error("Cannot Update New Department Employ Count", zap.Error(err))
			return err
		}

		if user.DepartmentID != nil && *user.DepartmentID != uint(*newDeploymentId) {
			var oldDepartment *department_models.Department
			if err := department_models.ValidScope(s.db.DB()).First(&oldDepartment, &user.DepartmentID).Error; err != nil {
				s.logger.Error("Cannot Find User's Department", zap.Error(err))
				return err
			}
			if err := oldDepartment.UpdateEmployCount(s.db.DB(), -1); err != nil {
				s.logger.Error("Cannot Update Old Department Employ Count", zap.Error(err))
				return err
			}
		}
	}

	return nil
}
