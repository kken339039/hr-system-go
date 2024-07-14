package controllers

import (
	"bytes"
	"encoding/json"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/department/dtos"
	"hr-system-go/internal/department/models"
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
	RunSpecs(t, "Department Controller Suite")
}

var (
	departmentController  *DepartmentController
	mockDepartmentService *mock_services.MockDepartmentService
	mockAuthService       *mock_services.MockAuthService
	router                *gin.Engine
	mockEnv               *env.Env
	mockLogger            *logger.Logger
)

var _ = Describe("DepartmentController", func() {
	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockEnv := env.NewEnv()
		mockLogger = logger.NewLogger(mockEnv)
		mockDepartmentService = &mock_services.MockDepartmentService{}
		mockAuthService = &mock_services.MockAuthService{}
		departmentController = NewDepartmentController(mockLogger, mockDepartmentService, mockAuthService)
		router = gin.Default()
		departmentController.RegisterRoutes(router)
	})
	Describe("listDepartments", func() {
		It("should return a list of departments", func() {
			departments := []models.Department{
				{Name: "HR"},
				{Name: "IT"},
			}
			totalRows := int64(2)

			mockDepartmentService.On("FindDepartments", mock.AnythingOfType("*utils.Pagination")).Return(departments, totalRows, nil)

			req, _ := http.NewRequest("GET", "/api/department", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Items"]).To(HaveLen(2))
		})
	})

	Describe("GetDepartment", func() {
		It("should return a specific department", func() {
			departmentID := 1
			department := &models.Department{Name: "HR"}
			department.ID = uint(departmentID)

			mockDepartmentService.On("FindDepartmentByID", departmentID).Return(department, nil)

			req, _ := http.NewRequest("GET", "/api/department/"+strconv.Itoa(departmentID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Name"]).To(Equal(department.Name))
		})
	})

	Describe("CreateDepartment", func() {
		It("should create a new department", func() {
			payload := dtos.CreateDepartmentRequest{Name: "New Department"}
			newDepartment := &models.Department{Name: "New Department"}
			newDepartment.ID = 3

			mockDepartmentService.On("CreateDepartment", payload).Return(newDepartment, nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/api/department", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Name"]).To(Equal(newDepartment.Name))
		})
	})

	Describe("UpdateDepartment", func() {
		It("should update an existing department", func() {
			departmentID := 1
			updateName := "Updated Department"
			payload := dtos.UpdateDepartmentRequest{Name: &updateName}
			updatedDepartment := &models.Department{Name: updateName}
			updatedDepartment.ID = uint(departmentID)

			mockDepartmentService.On("UpdateDepartmentByID", departmentID, payload).Return(updatedDepartment, nil)

			jsonPayload, _ := json.Marshal(payload)
			req, _ := http.NewRequest("PUT", "/api/department/"+strconv.Itoa(departmentID), bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			Expect(response["Name"]).To(Equal(updatedDepartment.Name))
		})
	})

	Describe("DeleteDepartment", func() {
		It("should delete a department", func() {
			departmentID := 1

			mockDepartmentService.On("DeleteDepartmentByID", departmentID).Return(nil)

			req, _ := http.NewRequest("DELETE", "/api/department/"+strconv.Itoa(departmentID), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNoContent))
		})
	})
})
