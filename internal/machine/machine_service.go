package machine

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.MachineResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MachineResponse]
	FindById(ginContext *gin.Context, machineId uint64) *model.MachineResponse
	FindByFacilityId(ginContext *gin.Context, facilityId uint64) []*model.MachineResponse
	Create(ginContext *gin.Context, createMachineRequest *model.CreateMachineRequest)
	Update(ginContext *gin.Context, updateMachineRequest *model.UpdateMachineRequest)
	Delete(ginContext *gin.Context, deleteMachineRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.MachineResponse]
}
