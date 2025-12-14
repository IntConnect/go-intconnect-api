package node

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createNodeRequest *model.CreateNodeRequest)
	FindAll() []*model.NodeResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.NodeResponse]
	Update(ginContext *gin.Context, updateNodeRequest *model.UpdateNodeRequest)
	Delete(ginContext *gin.Context, deleteNodeRequest *model.DeleteNodeRequest)
}
