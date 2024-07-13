package services

import (
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/internal/department/dtos"
	"hr-system-go/internal/department/models"
	"hr-system-go/utils"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDepartmentService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DepartmentService Suite")
}

var (
	departmentService *DepartmentService
	mockEnv           *env.Env
	mockLogger        *logger.Logger
	mockDB            *mysql.MySqlStore
)

var _ = BeforeSuite(func() {
	mockEnv = env.NewEnv()
	mockLogger = logger.NewLogger(mockEnv)
	mockDB = mysql.NewMySqlStore(mockEnv, mockLogger)
	departmentService = NewDepartmentService(mockLogger, mockDB)

	mockDB.Connect(
		mockEnv.GetEnv("DB_USER"),
		mockEnv.GetEnv("DB_PASSWORD"),
		mockEnv.GetEnv("DB_DATABASE"),
		mockEnv.GetEnv("DB_HOST"),
		mockEnv.GetEnv("DB_PORT"),
		mockEnv.GetEnv("DB_PARAMS"),
	)

	mockDB.DB().AutoMigrate(&models.Department{})
})

var _ = AfterSuite(func() {
	mockDB.DB().Migrator().DropTable(&models.Department{})
	mockDB.Close()
})

var _ = Describe("DepartmentService", func() {
	Describe("FindDepartments", func() {
		BeforeEach(func() {
			_ = mockDB.DB().Exec("truncate table department").Error
		})

		It("should return departments and total count", func() {
			pagination := utils.Pagination{Page: 1, Limit: 10, Sort: "id asc"}
			departments := []*models.Department{{Name: "HR"}, {Name: "IT"}}

			mockDB.DB().Create(departments)

			result, totalCount, err := departmentService.FindDepartments(pagination)

			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(2))
			Expect(totalCount).To(Equal(int64(2)))
		})
	})

	Describe("FindDepartmentByID", func() {
		It("should return a department when found", func() {
			department := &models.Department{
				Name: "HR",
			}
			mockDB.DB().Create(&department)
			result, err := departmentService.FindDepartmentByID(int(department.ID))

			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal(department.ID))
			Expect(result.Name).To(Equal("HR"))
		})

		It("should return an error when department is not found", func() {
			departmentID := 99

			result, err := departmentService.FindDepartmentByID(departmentID)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})

	Describe("CreateDepartment", func() {
		It("should create a new department", func() {
			payload := dtos.CreateDepartmentRequest{Name: "New Dept", Descriptions: new(string)}
			*payload.Descriptions = "New department description"

			result, err := departmentService.CreateDepartment(payload)

			Expect(err).To(BeNil())
			Expect(result.Name).To(Equal(payload.Name))
			Expect(result.Descriptions).To(Equal(*payload.Descriptions))
		})
	})

	Describe("UpdateDepartmentByID", func() {
		It("should update an existing department", func() {
			requestName := "Updated Dept"
			payload := dtos.UpdateDepartmentRequest{Name: &requestName}

			department := &models.Department{
				Name: requestName,
			}
			mockDB.DB().Create(&department)
			result, err := departmentService.UpdateDepartmentByID(int(department.ID), payload)

			Expect(err).To(BeNil())
			Expect(result.ID).To(Equal(department.ID))
			Expect(result.Name).To(Equal(requestName))
		})
	})

	Describe("DeleteDepartmentByID", func() {
		It("should mark a department as removed", func() {
			departmentName := "Deleted Dept"
			department := &models.Department{
				Name: departmentName,
			}

			mockDB.DB().Create(&department)

			err := departmentService.DeleteDepartmentByID(int(department.ID))

			Expect(err).To(BeNil())
		})
	})
})
