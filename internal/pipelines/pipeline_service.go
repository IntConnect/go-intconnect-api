package pipeline

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createNodeDto *model.CreateNodeDto)
	FindAll() []*model.NodeResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.NodeResponse]
	Update(ginContext *gin.Context, updateNodeDto *model.UpdateNodeDto)
	Delete(ginContext *gin.Context, deleteNodeDto *model.DeleteNodeDto)
}
