package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/attendance/dtos"
	"hr-system-go/internal/attendance/models"
	"hr-system-go/internal/auth/constants"
	user_models "hr-system-go/internal/user/models"
	mock_services "hr-system-go/mocks/services"
	"hr-system-go/utils"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func TestAttendanceController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Leave & ClockRecord Controller Suite")
}

var (
	leaveController        *LeaveController
	clockRecordController  *ClockRecordController
	mockLeaveService       *mock_services.MockLeaveService
	mockClockRecordService *mock_services.MockClockRecordService
	mockAuthService        *mock_services.MockAuthService
	router                 *gin.Engine
	mockEnv                *env.Env
	mockLogger             *logger.Logger
)

var _ = Describe("LeavesController and ClockRecordController", func() {
	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockEnv := env.NewEnv()
		mockLogger = logger.NewLogger(mockEnv)
		mockLeaveService = &mock_services.MockLeaveService{}
		mockAuthService = &mock_services.MockAuthService{}
		mockClockRecordService = &mock_services.MockClockRecordService{}
		leaveController = NewLeaveController(mockLogger, mockLeaveService, mockAuthService)
		clockRecordController = NewClockRecordController(mockLogger, mockClockRecordService, mockAuthService)
		router = gin.Default()
		leaveController.RegisterRoutes(router)
		clockRecordController.RegisterRoutes(router)
	})

	Describe("LeaveController", func() {
		Describe("listLeaves", func() {
			It("should return a list of leaves", func() {
				userId := "1"
				userID, _ := strconv.Atoi(userId)
				leaves := []models.Leave{
					{UserID: uint(userID), StartDate: time.Now(), EndDate: time.Now().Add(24 * time.Hour)},
					{UserID: uint(userID), StartDate: time.Now().Add(48 * time.Hour), EndDate: time.Now().Add(72 * time.Hour)},
				}
				totalRows := int64(2)

				mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_LEAVE).Return(true)
				mockLeaveService.On("FindLeavesByUserID", userID, mock.AnythingOfType("*utils.Pagination")).Return(leaves, totalRows, nil)

				req, _ := http.NewRequest("GET", "/api/users/"+userId+"/leave", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				Expect(response["Items"]).To(HaveLen(2))
			})

			It("should return error when user is not authorized", func() {
				userId := "1"
				userID, _ := strconv.Atoi(userId)

				mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_LEAVE).Return(false)

				req, _ := http.NewRequest("GET", "/api/users/"+userId+"/leave", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusForbidden))
			})

			Describe("getLeave", func() {
				It("should return a specific leave", func() {
					userId := "1"
					leaveId := "1"
					userID, _ := strconv.Atoi(userId)
					leaveID, _ := strconv.Atoi(leaveId)

					leave := &models.Leave{UserID: uint(userID), StartDate: time.Now(), EndDate: time.Now().Add(24 * time.Hour)}

					mockUser := &user_models.User{}
					mockUser.ID = uint(userID)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(mockUser)
					mockLeaveService.On("FindLeaveByID", leaveID).Return(leave, nil)

					req, _ := http.NewRequest("GET", "/api/users/"+userId+"/leave/"+leaveId, nil)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusOK))

					var response map[string]interface{}
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(response["leave"]).NotTo(BeNil())
				})

				It("should return error when leave doesn't belong to user", func() {
					userId := "1"
					leaveId := "1"

					mockUser := &user_models.User{}
					mockUser.ID = uint(2)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(mockUser)

					req, _ := http.NewRequest("GET", "/api/users/"+userId+"/leave/"+leaveId, nil)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusForbidden))
				})
			})

			Describe("createLeave", func() {
				It("should create a new leave", func() {
					userId := "1"
					userID, _ := strconv.Atoi(userId)
					sts := "2024-07-01T15:04:05+08:00"
					ets := "2024-07-02T15:04:05+08:00"
					st, _ := utils.ParseDateTime(sts)
					et, _ := utils.ParseDateTime(ets)
					leaveType := "annual"
					payload := dtos.CreateLeaveRequest{
						StartDate: &sts,
						EndDate:   &ets,
						LeaveType: &leaveType,
					}

					leave := &models.Leave{UserID: uint(userID), StartDate: st, EndDate: et, LeaveType: leaveType}
					user := &user_models.User{}
					user.ID = uint(userID)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)
					mockLeaveService.On("CreateLeaveByUser", mock.AnythingOfType("*models.User"), payload).Return(leave, nil)

					jsonPayload, _ := json.Marshal(payload)
					req, _ := http.NewRequest("POST", "/api/users/"+userId+"/leave", bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusOK))

					var response map[string]interface{}
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(response["leave"]).NotTo(BeNil())
				})

				It("should return error when creating leave for another user", func() {
					userId := "1"
					userID, _ := strconv.Atoi(userId)
					sts := "2024-07-01T15:04:05+08:00"
					ets := "2024-07-02T15:04:05+08:00"
					leaveType := "annual"
					payload := dtos.CreateLeaveRequest{
						StartDate: &sts,
						EndDate:   &ets,
						LeaveType: &leaveType,
					}
					user := &user_models.User{}
					user.ID = uint(userID + 1)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)

					jsonPayload, _ := json.Marshal(payload)
					req, _ := http.NewRequest("POST", "/api/users/"+userId+"/leave", bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusForbidden))
				})
			})

			Describe("updateLeave", func() {
				It("should update an existing leave", func() {
					userId := "1"
					leaveId := "1"
					userID, _ := strconv.Atoi(userId)
					leaveID, _ := strconv.Atoi(leaveId)

					sts := "2024-07-01T15:04:05+08:00"
					ets := "2024-07-02T15:04:05+08:00"
					st, _ := utils.ParseDateTime(sts)
					et, _ := utils.ParseDateTime(ets)
					leaveType := "sick"
					payload := dtos.UpdateLeaveRequest{
						StartDate: &sts,
						EndDate:   &ets,
						LeaveType: &leaveType,
					}

					updatedLeave := &models.Leave{UserID: uint(userID), StartDate: st, EndDate: et, LeaveType: *payload.LeaveType}
					updatedLeave.ID = uint(leaveID)
					user := &user_models.User{}
					user.ID = uint(userID)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)
					mockLeaveService.On("UpdateLeaveByID", leaveID, payload).Return(updatedLeave, nil)

					jsonPayload, _ := json.Marshal(payload)
					req, _ := http.NewRequest("PUT", "/api/users/"+userId+"/leave/"+leaveId, bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusOK))

					var response map[string]interface{}
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(response["leave"]).NotTo(BeNil())
				})

				It("should return error when updating leave for another user", func() {
					userId := "1"
					leaveId := "1"
					userID, _ := strconv.Atoi(userId)

					sts := "2024-07-01T15:04:05+08:00"
					ets := "2024-07-02T15:04:05+08:00"
					leaveType := "sick"
					payload := dtos.UpdateLeaveRequest{
						StartDate: &sts,
						EndDate:   &ets,
						LeaveType: &leaveType,
					}
					user := &user_models.User{}
					user.ID = uint(userID + 1)

					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)

					jsonPayload, _ := json.Marshal(payload)
					req, _ := http.NewRequest("PUT", "/api/users/"+userId+"/leave/"+leaveId, bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusForbidden))
				})
			})

			Describe("deleteLeave", func() {
				It("should delete a leave", func() {
					userId := "1"
					leaveId := "1"
					userID, _ := strconv.Atoi(userId)
					leaveID, _ := strconv.Atoi(leaveId)

					user := &user_models.User{}
					user.ID = uint(userID)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)
					mockLeaveService.On("DeleteLeaveByID", leaveID).Return(nil)

					req, _ := http.NewRequest("DELETE", "/api/users/"+userId+"/leave/"+leaveId, nil)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should return error when deleting leave for another user", func() {
					userId := "1"
					leaveId := "1"
					userID, _ := strconv.Atoi(userId)

					user := &user_models.User{}
					user.ID = uint(userID + 1)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)

					req, _ := http.NewRequest("DELETE", "/api/users/"+userId+"/leave/"+leaveId, nil)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusForbidden))
				})

				It("should return error when leave deletion fails", func() {
					userId := "1"
					leaveId := "1"
					userID, _ := strconv.Atoi(userId)
					leaveID, _ := strconv.Atoi(leaveId)

					user := &user_models.User{}
					user.ID = uint(userID)
					mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)
					mockLeaveService.On("DeleteLeaveByID", leaveID).Return(errors.New("deletion failed"))

					req, _ := http.NewRequest("DELETE", "/api/users/"+userId+"/leave/"+leaveId, nil)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
	Describe("ClockRecordController", func() {
		Describe("listClockRecord", func() {
			It("should return a list of clock records", func() {
				userId := "1"
				userID, _ := strconv.Atoi(userId)
				records := []models.ClockRecord{
					{UserID: uint(userID)},
					{UserID: uint(userID)},
				}
				totalRows := int64(2)

				mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_CLOCK_RECORD).Return(true)
				mockClockRecordService.On("FindClockRecordsByUserID", userID, mock.AnythingOfType("*utils.Pagination")).Return(records, totalRows, nil)

				req, _ := http.NewRequest("GET", "/api/users/"+userId+"/clockRecord", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				Expect(response["Items"]).To(HaveLen(2))
			})

			It("should return error when user is not authorized", func() {
				userId := "1"
				userID, _ := strconv.Atoi(userId)

				mockAuthService.On("AbleToAccessOtherUserData", mock.Anything, userID, constants.ABILITY_ALL_GRANTS_CLOCK_RECORD).Return(false)

				req, _ := http.NewRequest("GET", "/api/users/"+userId+"/clockRecord", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusForbidden))
			})
		})

		Describe("touchClockRecord", func() {
			It("should create a new clock record", func() {
				userId := "1"
				userID, _ := strconv.Atoi(userId)
				user := &user_models.User{}
				user.ID = uint(userID)
				record := &models.ClockRecord{UserID: uint(userID)}
				record.ID = uint(1)

				mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)
				mockClockRecordService.On("ClockByUser", user).Return(record, nil)

				req, _ := http.NewRequest("POST", "/api/users/"+userId+"/clockRecord/clock", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				Expect(response["leave"]).NotTo(BeNil())
			})

			It("should return error when creating clock record for another user", func() {
				userId := "1"
				userID, _ := strconv.Atoi(userId)
				user := &user_models.User{}
				user.ID = uint(userID + 1)

				mockAuthService.On("GetCurrentUser", mock.Anything).Return(user)

				req, _ := http.NewRequest("POST", "/api/users/"+userId+"/clockRecord/clock", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusForbidden))
			})
		})
	})
})
