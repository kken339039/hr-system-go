package services

import (
	"hr-system-go/internal/user/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) AuthTokenWrapper(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c)
	}
}

func (m *MockAuthService) AuthUserAbilityWrapper(handler gin.HandlerFunc, ability string) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c)
	}
}

func (m *MockAuthService) AbleToAccessOtherUserData(ctx *gin.Context, userID int, ability string) bool {
	args := m.Called(ctx, userID, ability)
	return args.Bool(0)
}

func (m *MockAuthService) GetCurrentUser(ctx *gin.Context) *models.User {
	args := m.Called(ctx)
	return args.Get(0).(*models.User)
}

func (m *MockAuthService) GenerateToken(userID uint, username string) (string, error) {
	args := m.Called(userID, username)
	return args.String(0), args.Error(1)
}
