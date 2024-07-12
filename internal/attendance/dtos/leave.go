package dtos

import (
	"hr-system-go/internal/attendance/models"
	"hr-system-go/utils"
	"time"
)

type LeaveListResponse struct {
	Items      []*LeaveResponse
	Pagination utils.PaginationResult
}

type LeaveResponse struct {
	Id        uint
	UserName  string
	StartDate time.Time
	EndDate   time.Time
	LeaveType string
	Status    string
}

type CreateLeaveRequest struct {
	StartDate *string `json:"startDate,omitempty"`
	EndDate   *string `json:"endDate,omitempty"`
	LeaveType *string `json:"leaveType,omitempty"`
}

type UpdateLeaveRequest struct {
	StartDate *string `json:"startDate,omitempty"`
	EndDate   *string `json:"endDate,omitempty"`
	LeaveType *string `json:"leaveType,omitempty"`
	Status    *string `json:"status,omitempty"`
}

func NewLeaveListResponse(leaves []models.Leave, totalRows int64, pagination utils.Pagination) *LeaveListResponse {
	items := []*LeaveResponse{}
	for _, leave := range leaves {
		items = append(items, NewLeaveResponse(&leave))
	}

	return &LeaveListResponse{
		Items: items,
		Pagination: utils.PaginationResult{
			Limit: pagination.Limit,
			Page:  pagination.Page,
			Total: totalRows,
			Sort:  pagination.Sort,
		},
	}
}

func NewLeaveResponse(leave *models.Leave) *LeaveResponse {
	res := &LeaveResponse{
		Id:        leave.ID,
		UserName:  leave.User.Name,
		StartDate: leave.StartDate,
		EndDate:   leave.EndDate,
		LeaveType: leave.LeaveType,
		Status:    leave.Status,
	}

	return res
}
