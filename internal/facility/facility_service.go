package facility

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.FacilityResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.FacilityResponse]
	FindById(facilityId uint64) *model.FacilityResponse
	Create(ginContext *gin.Context, createFacilityRequest *model.CreateFacilityRequest)
	Update(ginContext *gin.Context, updateFacilityRequest *model.UpdateFacilityRequest)
	Delete(ginContext *gin.Context, deleteFacilityRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.FacilityResponse]
}
