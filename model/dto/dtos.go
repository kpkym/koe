package dto

import "github.com/kpkym/koe/model/domain"

type SearchResponse struct {
	Pagination `json:"pagination"`
	Works      []domain.WorkDomain `json:"works"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage,omitempty"`
	PageSize    int `json:"pageSize,omitempty"`
	TotalCount  int `json:"totalCount,omitempty"`
}
