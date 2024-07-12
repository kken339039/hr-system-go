package dtos

import (
	"hr-system-go/internal/user/models"
	"hr-system-go/utils"
)

type UserListResponse struct {
	Items      []*UserResponse
	Pagination utils.PaginationResult
}

type UserResponse struct {
	Name           *string
	Email          *string
	Age            *int
	Status         *string
	Salary         *float64
	RoleName       *string
	DepartmentName *string
}

type UpdateUserRequest struct {
	Name         *string  `json:"name,omitempty"`
	Email        *string  `json:"email,omitempty"`
	Age          *int     `json:"age,omitempty"`
	Status       *string  `json:"status,omitempty"`
	Salary       *float64 `json:"salary,omitempty"`
	RoleID       *int     `json:"roleId,omitempty"`
	DepartmentID *int     `json:"departmentId,omitempty"`
}

func NewUserListResponse(users []models.User, totalRows int64, pagination utils.Pagination) *UserListResponse {
	items := []*UserResponse{}
	for _, user := range users {
		items = append(items, NewUserResponse(&user))
	}

	return &UserListResponse{
		Items: items,
		Pagination: utils.PaginationResult{
			Limit: pagination.Limit,
			Page:  pagination.Page,
			Total: totalRows,
			Sort:  pagination.Sort,
		},
	}
}

func NewUserResponse(user *models.User) *UserResponse {
	res := &UserResponse{
		Name:   &user.Name,
		Email:  &user.Email,
		Age:    &user.Age,
		Status: &user.Status,
		Salary: user.Salary,
	}

	if user.Role != nil {
		res.RoleName = &user.Role.Name
	}
	if user.Department != nil {
		res.DepartmentName = &user.Department.Name
	}

	return res
}
