package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit int
	Page  int
	Sort  string
}

type PaginationResult struct {
	Limit int
	Page  int
	Total int64
	Sort  string
}

func NewPagination(ctx *gin.Context) Pagination {
	pagination := Pagination{}
	if limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")); err == nil {
		pagination.Limit = limit
	}

	if page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")); err == nil {
		pagination.Page = page
	}

	pagination.Sort = ctx.DefaultQuery("sort", "id desc")
	return pagination
}

func (p Pagination) Offset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}
