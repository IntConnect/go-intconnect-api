package facility

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createFacilityRequest *model.CreateFacilityRequest)
	FindAll() []*model.FacilityResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.FacilityResponse]
	Update(ginContext *gin.Context, updateFacilityRequest *model.UpdateFacilityRequest)
	Delete(ginContext *gin.Context, deleteFacilityRequest *model.DeleteFacilityRequest)
}
