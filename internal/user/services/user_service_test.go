package services

import (
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	auth_models "hr-system-go/internal/auth/models"
	department_models "hr-system-go/internal/department/models"
	"hr-system-go/internal/user/dtos"
	"hr-system-go/internal/user/models"
	"hr-system-go/utils"
	"testing"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")
}

var (
	userService *UserService
	mockEnv     *env.Env
	mockLogger  *logger.Logger
	mockDB      *mysql.MySqlStore
)

var _ = BeforeSuite(func() {
	mockEnv = env.NewEnv()
	mockLogger = logger.NewLogger(mockEnv)
	mockDB = mysql.NewMySqlStore(mockEnv, mockLogger)
	userService = NewUserService(mockLogger, mockDB)

	mockDB.Connect(
		mockEnv.GetEnv("DB_USER"),
		mockEnv.GetEnv("DB_PASSWORD"),
		mockEnv.GetEnv("DB_DATABASE"),
		mockEnv.GetEnv("DB_HOST"),
		mockEnv.GetEnv("DB_PORT"),
		mockEnv.GetEnv("DB_PARAMS"),
	)

	mockDB.DB().AutoMigrate(&models.User{}, &auth_models.Role{}, &department_models.Department{})
})

var _ = AfterSuite(func() {
	mockDB.DB().Migrator().DropTable(&models.User{}, &auth_models.Role{}, &department_models.Department{})
	mockDB.Close()
})

var _ = Describe("UserService", func() {
	Describe("RegisterUser", func() {
		It("should register a new user successfully", func() {
			user := &models.User{
				Name:  "John Doe",
				Email: faker.Email(),
			}
			password := "password123"

			err := userService.RegisterUser(user, password)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.PasswordEncrypt).ShouldNot(BeEmpty())
			Expect(user.JoinDate).ShouldNot(BeZero())
		})
	})

	Describe("FindUsers", func() {
		BeforeEach(func() {
			_ = mockDB.DB().Exec("truncate table user").Error
		})

		It("should find users with pagination", func() {
			pagination := utils.Pagination{
				Page:  1,
				Limit: 10,
				Sort:  "id asc",
			}

			mockUsers := []*models.User{
				{
					Name:  "John",
					Email: faker.Email(),
				},
				{
					Name:  "Marry",
					Email: faker.Email(),
				},
			}
			mockDB.DB().Create(mockUsers)

			users, totalCount, err := userService.FindUsers(pagination)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(users).To(HaveLen(2))
			Expect(totalCount).To(Equal(int64(2)))
		})
	})

	Describe("FindUserByEmail", func() {
		It("should find a user by email", func() {
			mockUser := &models.User{
				Email: faker.Email(),
			}
			mockDB.DB().Create(mockUser)
			user, err := userService.FindUserByEmail(mockUser.Email)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.Email).To(Equal(mockUser.Email))
		})
	})

	Describe("FindUserByID", func() {
		It("should find a user by ID", func() {
			mockUser := &models.User{
				Name:  "John",
				Email: faker.Email(),
			}

			mockDB.DB().Create(&mockUser)

			user, err := userService.FindUserByID(int(mockUser.ID))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.ID).To(Equal(user.ID))
		})
	})

	Describe("UpdateUserByID", func() {
		It("should update a user by ID", func() {
			mockDepartment := &department_models.Department{
				Name: "MOMO",
			}
			mockDB.DB().Create(&mockDepartment)
			departmentID := int(mockDepartment.ID)

			updatedName := "Updated John"
			payload := dtos.UpdateUserRequest{
				Name:         &updatedName,
				DepartmentID: &departmentID,
			}

			mockUser := &models.User{
				Name:  "John",
				Email: faker.Email(),
			}
			mockDB.DB().Create(&mockUser)

			mockDB.DB().Create(&mockDepartment)
			user, err := userService.UpdateUserByID(int(mockUser.ID), payload)
			var department *department_models.Department
			mockDB.DB().First(&department, departmentID)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.Name).To(Equal("Updated John"))
			Expect(department.EmployCount).To(Equal(1))
		})
	})

	Describe("DeleteUserByID", func() {
		It("should delete a user by ID", func() {
			mockDepartment := &department_models.Department{
				Name: "MOMO",
			}
			mockDB.DB().Create(&mockDepartment)
			mockUser := &models.User{
				Name:         "JohnM",
				Email:        faker.Email(),
				DepartmentID: &mockDepartment.ID,
			}
			mockDB.DB().Create(&mockUser)

			err := userService.DeleteUserByID(int(mockUser.ID))

			var department *department_models.Department
			mockDB.DB().First(&department, mockDepartment.ID)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(department.EmployCount).To(Equal(0))
		})
	})

	Describe("UpdatePassword", func() {
		It("should update user's password", func() {
			mockUser := &models.User{
				Email:           faker.Email(),
				PasswordEncrypt: "$2a$10$abcdefghijklmnopqrstuvwxyz012345",
			}
			mockDB.DB().Create(&mockUser)
			newPassword := "newpassword123"

			err := userService.UpdatePassword(mockUser, newPassword)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
