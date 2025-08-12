package model

type PaginationRequest struct {
	Page        int    `form:"page" binding:"min=1"`
	Size        int    `form:"size" binding:"min=1,max=100"`
	Sort        string `form:"sortBy"`
	Order       string `form:"sortDir" binding:"omitempty,oneof=asc desc"`
	SearchQuery string `form:"searchQuery" binding:"omitempty"`
	StartDate   string `form:"startDate" binding:"omitempty"`
	EndDate     string `form:"endDate" binding:"omitempty"`
	Filter      string `form:"filter"`
}
type PaginationResponse[T any] struct {
	Data        []T   `json:"data"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
	CurrentPage int   `json:"currentPage"`
}

func NewPaginationRequest() PaginationRequest {
	return PaginationRequest{
		Page:  1,
		Size:  10,
		Sort:  "id",
		Order: "asc",
	}
}

func PaginationResponseFactory[T any](paginationReq *PaginationRequest) (PaginationResponse[*T], int, string) {
	paginationResp := PaginationResponse[*T]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	orderClause += " " + paginationReq.Order
	return paginationResp, offsetVal, orderClause
}
