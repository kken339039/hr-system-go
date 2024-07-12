package dtos

import (
	"hr-system-go/internal/attendance/models"
	"hr-system-go/utils"
	"time"
)

type ClockRecordListResponse struct {
	Items      []*ClockRecordResponse
	Pagination utils.PaginationResult
}

type ClockRecordResponse struct {
	Id       uint
	UserName string
	ClockIn  time.Time
	ClockOut *time.Time
}

type ClockInRequest struct {
	ClockIn *string `json:"clockIn,omitempty"`
}

type ClockOutRequest struct {
	ClockOut *string `json:"clockOut,omitempty"`
}

type UpdateClockRecordRequest struct {
	StartDate       *string `json:"startDate,omitempty"`
	EndDate         *string `json:"endDate,omitempty"`
	ClockRecordType *string `json:"ClockRecordType,omitempty"`
	Status          *string `json:"status,omitempty"`
}

func NewClockRecordListResponse(ClockRecords []models.ClockRecord, totalRows int64, pagination utils.Pagination) *ClockRecordListResponse {
	items := []*ClockRecordResponse{}
	for _, ClockRecord := range ClockRecords {
		items = append(items, NewClockRecordResponse(&ClockRecord))
	}

	return &ClockRecordListResponse{
		Items: items,
		Pagination: utils.PaginationResult{
			Limit: pagination.Limit,
			Page:  pagination.Page,
			Total: totalRows,
			Sort:  pagination.Sort,
		},
	}
}

func NewClockRecordResponse(clockRecord *models.ClockRecord) *ClockRecordResponse {
	res := &ClockRecordResponse{
		Id:       clockRecord.ID,
		UserName: clockRecord.User.Name,
		ClockIn:  clockRecord.ClockIn,
		ClockOut: clockRecord.ClockOut,
	}

	return res
}
