package dtos

import (
	"hr-system-go/internal/department/models"
	"hr-system-go/utils"
)

type DepartmentListResponse struct {
	Items      []*DepartmentResponse
	Pagination utils.PaginationResult
}

type DepartmentResponse struct {
	Id          uint
	Name        string
	Description string
	Status      string
	EmployCount int
}

type CreateDepartmentRequest struct {
	Name         string  `json:"name"`
	Descriptions *string `json:"descriptions,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name         *string `json:"name"`
	Descriptions *string `json:"descriptions,omitempty"`
	Status       *string `json:"status,omitempty"`
}

func NewDepartmentListResponse(departments []models.Department, totalRows int64, pagination utils.Pagination) *DepartmentListResponse {
	items := []*DepartmentResponse{}
	for _, department := range departments {
		items = append(items, NewDepartmentResponse(&department))
	}

	return &DepartmentListResponse{
		Items: items,
		Pagination: utils.PaginationResult{
			Limit: pagination.Limit,
			Page:  pagination.Page,
			Total: totalRows,
			Sort:  pagination.Sort,
		},
	}
}

func NewDepartmentResponse(department *models.Department) *DepartmentResponse {
	res := &DepartmentResponse{
		Id:          department.ID,
		Name:        department.Name,
		Description: department.Descriptions,
		Status:      department.Status,
		EmployCount: department.EmployCount,
	}

	return res
}
