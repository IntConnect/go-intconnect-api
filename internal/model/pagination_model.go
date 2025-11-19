package model

type PaginationRequest struct {
	Page        int    `form:"page,default=1"`
	Size        int    `form:"size,default=10"`
	Sort        string `form:"sort,default=id"`
	Order       string `form:"order,default=asc"`
	SearchQuery string `form:"query"`
}

type PaginationQuery struct {
	Offset      int
	Limit       int
	OrderClause string
	SearchQuery string
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
