package pipeline

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createPipelineRequest *model.CreatePipelineRequest)
	FindAll() []*model.PipelineResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.PipelineResponse]
	Update(ginContext *gin.Context, updatePipelineRequest *model.UpdatePipelineRequest)
	Delete(ginContext *gin.Context, deletePipelineRequest *model.DeletePipelineRequest)
}
