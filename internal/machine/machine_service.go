package machine

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.MachineResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MachineResponse]
	Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest) *model.PaginatedResponse[*model.MachineResponse]
	Update(ginContext *gin.Context, updateMachineRequest *model.UpdateMachineRequest) *model.PaginatedResponse[*model.MachineResponse]
	Delete(ginContext *gin.Context, deleteMachineRequest *model.DeleteMachineRequest) *model.PaginatedResponse[*model.MachineResponse]
}
