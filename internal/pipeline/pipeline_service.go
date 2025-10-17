package pipeline

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createPipelineDto *model.CreatePipelineDto)
	FindAll() []*model.PipelineResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.PipelineResponse]
	Update(ginContext *gin.Context, updatePipelineDto *model.UpdatePipelineDto)
	Delete(ginContext *gin.Context, deletePipelineDto *model.DeletePipelineDto)
}
