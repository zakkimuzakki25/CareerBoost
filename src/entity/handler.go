package entity

import "math"

type HTTPResponse struct {
	Message    string           `json:"message"`
	IsSuccess  bool             `json:"isSuccess"`
	Data       interface{}      `json:"data"`
	Pagination *PaginationParam `json:"pagination"`
}

type PaginationParam struct {
	Limit           int64 `json:"limit" form:"limit"`
	Page            int64 `json:"page" form:"page"`
	Offset          int64 `json:"offset"`
	TotalElements   int64 `json:"totalElements"`
	CurrentElements int64 `json:"currentElements"`
	TotalPages      int64 `json:"totalPages"`
	CurrentPages    int64 `json:"currentPages"`
}

func (pp *PaginationParam) FormatPagination() {
	if pp.Limit == 0 {
		pp.Limit = 5
	}

	if pp.Page == 0 {
		pp.Page = 1
	}

	pp.Offset = (pp.Page - 1) * pp.Limit
}

func (pp *PaginationParam) ProcessPagination(totalElements int64) {
	pp.TotalElements = totalElements
	pp.TotalPages = int64(math.Ceil(float64(pp.TotalElements) / float64(pp.Limit)))
	pp.CurrentPages = pp.Page

	if totalElements < pp.Limit {
		pp.CurrentElements = pp.TotalElements
	}
}

type Filter struct {
	InterestID []int  `form:"interestID[]"`
	Key        string `form:"search"`
}

// type mentorParam struct {
// 	ID uint `uri:"mentor_id" gorm:"column:id"`
// }
