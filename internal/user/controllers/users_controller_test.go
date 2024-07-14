package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/auth/constants"
	dtos "hr-system-go/internal/user/dtos"
	"hr-system-go/internal/user/models"
	mock_services "hr-system-go/mocks/services"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func TestDepartmentController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Controller Suite")
}

var (
	userController  *UsersController
	mockUserService *mock_services.MockUserService
	mockAuthService *mock_services.MockAuthService
	router          *gin.Engine
	mockEnv         *env.Env
	mockLogger      *logger.Logger
)

var _ = Describe("UserController", func() {
	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockEnv := env.NewEnv()
		mockLogger = logger.NewLogger(mockEnv)
		mockUserService = &mock_services.MockUserService{}
		mockAuthService = &mock_services.MockAuthService{}
		userController = NewUsersController(mockLogger, mockUserService, mockAuthService)
		router = gin.Default()
		userController.RegisterRoutes(router)
	})

	Describe("listUsers", func() {
		It("should return a list of users", func() {
			users := []models.User{
				{Name: "User1"},
				{Name: "User2"},
			}
			totalRows := int64(2)

			mockUserService.On("FindUsers", mock.AnythingOfType("*utils.Pagination")).Return(users, totalRows, nil)

			req, _ := http.NewRequest("GET", "/api/users", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Items"]).To(HaveLen(2))
		})
	})

	Describe("GetUser", func() {
		It("should return a specific user", func() {
			userID := 1
			user := &models.User{Name: "User1"}
			user.ID = uint(userID)

			mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_USER).Return(true)
			mockUserService.On("FindUserByID", userID).Return(user, nil)

			req, _ := http.NewRequest("GET", "/api/users/"+strconv.Itoa(userID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Name"]).To(Equal(user.Name))
		})

		It("should return forbidden when not able to access user data", func() {
			userID := 1

			mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_USER).Return(false)

			req, _ := http.NewRequest("GET", "/api/users/"+strconv.Itoa(userID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusForbidden))
		})
	})

	Describe("UpdateUser", func() {
		It("should update a user", func() {
			userID := 1
			updatedName := "Updated User"
			payload := dtos.UpdateUserRequest{Name: &updatedName}
			updatedUser := &models.User{Name: updatedName}
			updatedUser.ID = uint(userID)

			mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_USER).Return(true)
			mockUserService.On("UpdateUserByID", userID, payload).Return(updatedUser, nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("PUT", "/api/users/"+strconv.Itoa(userID), bytes.NewBuffer(jsonPayload))

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Name"]).To(Equal(updatedUser.Name))
		})
	})

	Describe("DeleteUser", func() {
		It("should delete a user", func() {
			userID := 1

			mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_USER).Return(true)
			mockUserService.On("DeleteUserByID", userID).Return(nil)

			req, _ := http.NewRequest("DELETE", "/api/users/"+strconv.Itoa(userID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNoContent))
		})

		It("should return an error when delete fails", func() {
			userID := 1

			mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_USER).Return(true)
			mockUserService.On("DeleteUserByID", userID).Return(errors.New("delete failed"))

			req, _ := http.NewRequest("DELETE", "/api/users/"+strconv.Itoa(userID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
