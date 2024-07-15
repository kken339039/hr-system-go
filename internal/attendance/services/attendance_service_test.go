package services

import (
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/internal/attendance/dtos"
	"hr-system-go/internal/attendance/models"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/utils"
	"testing"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAttendanceService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Leave & ClockRecord Service Suite")
}

var (
	clockRecordService ClockRecordServiceInterface
	leaveService       LeaveServiceInterface
	mockEnv            *env.Env
	mockLogger         *logger.Logger
	mockDB             *mysql.MySqlStore
)

var _ = BeforeSuite(func() {
	mockEnv = env.NewEnv()
	mockLogger = logger.NewLogger(mockEnv)
	mockDB = mysql.NewMySqlStore(mockEnv, mockLogger)
	leaveService = NewLeaveService(mockLogger, mockDB)
	clockRecordService = NewClockRecordService(mockLogger, mockDB)

	mockDB.Connect(
		mockEnv.GetEnv("DB_USER"),
		mockEnv.GetEnv("DB_PASSWORD"),
		mockEnv.GetEnv("DB_DATABASE"),
		mockEnv.GetEnv("DB_HOST"),
		mockEnv.GetEnv("DB_PORT"),
		mockEnv.GetEnv("DB_PARAMS"),
	)

	mockDB.DB().AutoMigrate(&models.Leave{}, &models.ClockRecord{}, &user_models.User{})
})

var _ = AfterSuite(func() {
	mockDB.DB().Migrator().DropTable(&models.Leave{}, &models.ClockRecord{}, &user_models.User{})
	mockDB.Close()
})

var _ = Describe("LeavesService and ClockRecordService", func() {
	Describe("ClockRecordService", func() {
		Describe("FindClockRecordsByUserID", func() {
			BeforeEach(func() {
				_ = mockDB.DB().Exec("truncate table clock_record").Error
			})

			It("should find clock records by user ID with pagination", func() {
				pagination := utils.Pagination{
					Page:  1,
					Limit: 10,
					Sort:  "id asc",
				}

				ci1, _ := utils.ParseDateTime("2024-07-12T15:04:05+08:00")
				co1, _ := utils.ParseDateTime("2024-07-14T15:04:05+08:00")
				ci2, _ := utils.ParseDateTime("2024-07-20T15:04:05+08:00")
				co2, _ := utils.ParseDateTime("2024-07-25T15:04:05+08:00")

				mockUser := &user_models.User{
					Email: faker.Email(),
				}
				mockDB.DB().Create(mockUser)

				mockClockRecord := []*models.ClockRecord{
					{
						ClockIn:  ci1,
						ClockOut: &co1,
						User:     *mockUser,
					},
					{
						ClockIn:  ci2,
						ClockOut: &co2,
						User:     *mockUser,
					},
				}
				mockDB.DB().Create(&mockClockRecord)

				records, totalCount, err := clockRecordService.FindClockRecordsByUserID(int(mockUser.ID), &pagination)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(records).To(HaveLen(2))
				Expect(totalCount).To(Equal(int64(2)))
			})
		})
		Describe("ClockByUser", func() {
			BeforeEach(func() {
				_ = mockDB.DB().Exec("truncate table clock_record").Error
			})
			Context("when clocking in", func() {
				It("should create a new clock record", func() {
					clockUser := &user_models.User{
						Email: faker.Email(),
					}
					mockDB.DB().Create(&clockUser)

					record, err := clockRecordService.ClockByUser(clockUser)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(record.UserID).To(Equal(clockUser.ID))
					Expect(record.ClockIn).ToNot(BeZero())
					Expect(record.ClockOut).To(BeZero())
				})
			})

			Context("when clocking out", func() {
				It("should update the existing clock record", func() {
					clockUser := &user_models.User{
						Email: faker.Email(),
					}
					mockDB.DB().Create(&clockUser)

					ci, _ := utils.ParseDateTime("2024-07-12T15:04:05+08:00")
					existedRecord := models.ClockRecord{
						User:    *clockUser,
						ClockIn: ci,
					}
					mockDB.DB().Create(&existedRecord)
					record, err := clockRecordService.ClockByUser(clockUser)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(record.UserID).To(Equal(clockUser.ID))
					Expect(record.ClockIn).ToNot(BeZero())
					Expect(record.ClockOut).ToNot(BeZero())
				})
			})
		})
	})
	Describe("LeavesService", func() {
		Describe("FindLeavesByUserID", func() {
			BeforeEach(func() {
				_ = mockDB.DB().Exec("truncate table leave").Error
			})

			It("should find leaves by user ID with pagination", func() {
				pagination := utils.Pagination{
					Page:  1,
					Limit: 10,
					Sort:  "id asc",
				}
				st1, _ := utils.ParseDateTime("2024-07-12T15:04:05+08:00")
				et1, _ := utils.ParseDateTime("2024-07-14T15:04:05+08:00")
				st2, _ := utils.ParseDateTime("2024-07-20T15:04:05+08:00")
				et2, _ := utils.ParseDateTime("2024-07-25T15:04:05+08:00")

				mockUser := &user_models.User{
					Email: faker.Email(),
				}
				mockDB.DB().Create(mockUser)

				mockLeaves := []*models.Leave{
					{
						StartDate: st1,
						EndDate:   et1,
						User:      *mockUser,
					},
					{
						StartDate: st2,
						EndDate:   et2,
						User:      *mockUser,
					},
				}
				mockDB.DB().Create(mockLeaves)

				leaves, totalCount, err := leaveService.FindLeavesByUserID(int(mockUser.ID), &pagination)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(leaves).To(HaveLen(2))
				Expect(totalCount).To(Equal(int64(2)))
			})
		})
		Describe("CreateLeaveByUser", func() {
			It("should create a new leave for a user", func() {
				startDate := "2024-07-01T15:04:05+08:00"
				endDate := "2024-07-02T15:04:05+08:00"
				leaveType := "annual"

				payload := dtos.CreateLeaveRequest{
					StartDate: &startDate,
					EndDate:   &endDate,
					LeaveType: &leaveType,
				}

				mockUser := &user_models.User{
					Email: faker.Email(),
				}
				mockDB.DB().Create(mockUser)

				leave, err := leaveService.CreateLeaveByUser(mockUser, payload)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(leave.UserID).To(Equal(mockUser.ID))
				Expect(leave.LeaveType).To(Equal(leaveType))
			})
		})
		Describe("UpdateLeaveByID", func() {
			It("should update a leave by ID", func() {
				newStatus := "approved"
				st1, _ := utils.ParseDateTime("2024-07-12T15:04:05+08:00")
				et1, _ := utils.ParseDateTime("2024-07-14T15:04:05+08:00")

				payload := dtos.UpdateLeaveRequest{
					Status: &newStatus,
				}

				mockUser := &user_models.User{
					Email: faker.Email(),
				}
				mockDB.DB().Create(mockUser)

				mockLeave := &models.Leave{
					StartDate: st1,
					EndDate:   et1,
					User:      *mockUser,
				}
				mockDB.DB().Create(mockLeave)

				leave, err := leaveService.UpdateLeaveByID(int(mockLeave.ID), payload)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(leave.ID).To(Equal(mockLeave.ID))
				Expect(leave.Status).To(Equal(newStatus))
			})
		})

		Describe("DeleteLeaveByID", func() {
			It("should delete a leave by ID", func() {
				st1, _ := utils.ParseDateTime("2024-07-12T15:04:05+08:00")
				et1, _ := utils.ParseDateTime("2024-07-14T15:04:05+08:00")
				mockUser := &user_models.User{
					Email: faker.Email(),
				}
				mockDB.DB().Create(mockUser)

				mockLeave := &models.Leave{
					StartDate: st1,
					EndDate:   et1,
					User:      *mockUser,
				}
				mockDB.DB().Create(mockLeave)

				err := leaveService.DeleteLeaveByID(int(mockLeave.ID))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Describe("FindLeaveByID", func() {
			It("should find a leave by ID", func() {
				st1, _ := utils.ParseDateTime("2024-07-12T15:04:05+08:00")
				et1, _ := utils.ParseDateTime("2024-07-14T15:04:05+08:00")
				mockUser := &user_models.User{
					Email: faker.Email(),
				}
				mockDB.DB().Create(mockUser)

				mockLeave := &models.Leave{
					StartDate: st1,
					EndDate:   et1,
					User:      *mockUser,
				}
				mockDB.DB().Create(mockLeave)

				leave, err := leaveService.FindLeaveByID(int(mockLeave.ID))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(leave.ID).To(Equal(mockLeave.ID))
			})
		})
	})
})
