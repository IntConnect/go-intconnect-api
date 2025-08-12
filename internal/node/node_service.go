package node

import (
	"github.com/gin-gonic/gin"
	"go-intconnect-api/internal/model"
)

type Service interface {
	Create(ginContext *gin.Context, createNodeDto *model.CreateNodeDto)
	FindAll() []*model.NodeResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.NodeResponse]
	Update(ginContext *gin.Context, updateNodeDto *model.UpdateNodeDto)
	Delete(ginContext *gin.Context, deleteNodeDto *model.DeleteNodeDto)
}
