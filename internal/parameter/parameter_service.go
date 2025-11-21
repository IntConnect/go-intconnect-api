package parameter

import (
	"go-intconnect-api/internal/model"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createParameterRequest *model.CreateParameterRequest, modelFile *multipart.FileHeader)
	FindAll() []*model.ParameterResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.ParameterResponse]
	Update(ginContext *gin.Context, updateParameterRequest *model.UpdateParameterRequest)
	Delete(ginContext *gin.Context, deleteParameterRequest *model.DeleteParameterRequest)
}
