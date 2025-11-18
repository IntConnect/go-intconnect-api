package permission

import (
	"go-intconnect-api/internal/model"
)

type Service interface {
	FindAll() []*model.PermissionResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.PermissionResponse]
}
