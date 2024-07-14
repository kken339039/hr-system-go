package controllers

import (
	"bytes"
	"encoding/json"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	user_models "hr-system-go/internal/user/models"
	mock_services "hr-system-go/mocks/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSessionController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Session Controller Suite")
}

var (
	sessionController *SessionsController
	mockUserService   *mock_services.MockUserService
	mockAuthService   *mock_services.MockAuthService
	router            *gin.Engine
	mockEnv           *env.Env
	mockLogger        *logger.Logger
)

var _ = Describe("UserController", func() {
	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockEnv := env.NewEnv()
		mockLogger = logger.NewLogger(mockEnv)
		mockUserService = &mock_services.MockUserService{}
		mockAuthService = &mock_services.MockAuthService{}
		sessionController = NewSessionsController(mockLogger, mockUserService, mockAuthService)
		router = gin.Default()
		sessionController.RegisterRoutes(router)
	})
	Describe("SignUp", func() {
		It("should register a new user and return a token", func() {
			payload := sessionBody{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			}
			user := &user_models.User{
				Name:  payload.Name,
				Email: payload.Email,
			}

			mockUserService.On("RegisterUser", user, payload.Password).Return(nil)
			mockAuthService.On("GenerateToken", mock.AnythingOfType("uint"), payload.Name).Return("token123", nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["token"]).To(Equal("token123"))
		})
	})

	Describe("SignIn", func() {
		It("should authenticate a user and return a token", func() {
			payload := sessionBody{
				Email:    "john@example.com",
				Password: "password123",
			}

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
			user := &user_models.User{Name: "John Doe", Email: "john@example.com", PasswordEncrypt: string(hashedPassword)}
			user.ID = uint(1)
			mockUserService.On("FindUserByEmail", payload.Email).Return(user, nil)
			mockAuthService.On("GenerateToken", user.ID, user.Name).Return("token123", nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["token"]).To(Equal("token123"))
		})
	})

	Describe("PasswordResetRequest", func() {
		It("should send a password reset request and return a token", func() {
			payload := passwordResetRequestBody{
				Email: "john@example.com",
			}

			user := &user_models.User{Name: "John Doe", Email: "john@example.com"}
			user.ID = uint(1)

			mockUserService.On("FindUserByEmail", payload.Email).Return(user, nil)
			mockAuthService.On("GenerateToken", user.ID, user.Name).Return("token123", nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/api/passwordResetRequest", bytes.NewBuffer(jsonPayload))
			// req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["token"]).To(Equal("token123"))
		})
	})

	Describe("ResetPassword", func() {
		It("should reset the user's password", func() {
			payload := resetPasswordBody{
				Email:       "john@example.com",
				NewPassword: "newpassword123",
			}

			user := &user_models.User{Name: "John Doe", Email: "john@example.com"}
			user.ID = uint(1)
			mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)
			mockUserService.On("UpdatePassword", user, payload.NewPassword).Return(nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/api/resetPassword", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNoContent))
		})
	})
})
