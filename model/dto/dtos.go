package dto

import "github.com/kpkym/koe/model/domain"

type PageRequest struct {
	Order string `form:"order"`
	Sort  string `form:"sort"`
	Page  int    `form:"page"`
	Size  int    `form:"size"`
}

type SearchResponse struct {
	Pagination `json:"pagination"`
	Works      []domain.WorkDomain `json:"works"`
}

type Pagination struct {
	CurrentPage int   `json:"currentPage,omitempty"`
	PageSize    int   `json:"pageSize,omitempty"`
	TotalCount  int64 `json:"totalCount,omitempty"`
}
