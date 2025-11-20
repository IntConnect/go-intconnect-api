package parameter

import (
	"go-intconnect-api/internal/model"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest, modelFile *multipart.FileHeader)
	FindAll() []*model.MachineResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.MachineResponse]
	Update(ginContext *gin.Context, updateMachineRequest *model.UpdateMachineRequest)
	Delete(ginContext *gin.Context, deleteMachineRequest *model.DeleteMachineRequest)
}
