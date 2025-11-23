package helper

import "go-intconnect-api/internal/model"

func WriteSuccess(message string, data interface{}) model.ResponseContractModel {
	return model.ResponseContractModel{
		Status:  true,
		Message: message,
		Data:    &data,
	}
}

func NewSuccessResponse[T any](message string, data T) *model.SuccessResponse[T] {
	return &model.SuccessResponse[T]{
		BaseResponse: model.BaseResponse{
			Success: true,
			Message: message,
		},
		Data: data,
	}
}

func NewErrorResponse(message string, errorDetail *model.ErrorDetail) *model.ErrorResponse {
	return &model.ErrorResponse{
		BaseResponse: model.BaseResponse{
			Success: false,
			Message: message,
		},
		Error: errorDetail,
	}
}

func NewPaginatedResponse[T any](
	message string,
	data []T,
	currentPage, pageSize int,
	totalItems int64,
) *model.PaginatedResponse[T] {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize != 0 {
		totalPages++
	}

	return &model.PaginatedResponse[T]{
		BaseResponse: model.BaseResponse{
			Success: true,
			Message: message,
		},
		Data: data,
		Pagination: &model.PaginationMeta{
			CurrentPage: currentPage,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalItems:  totalItems,
		},
	}
}

// NewPaginatedResponseFromResult helper untuk membuat paginated response dari result query
func NewPaginatedResponseFromResult[T any](
	message string,
	data []T,
	paginationReq *model.PaginationRequest,
	totalItems int64,
) *model.PaginatedResponse[T] {
	return NewPaginatedResponse(
		message,
		data,
		paginationReq.Page,
		paginationReq.Size,
		totalItems,
	)
}
