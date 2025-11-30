package facility

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	facilityService Service
	viperConfig     *viper.Viper
}

func NewHandler(facilityService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		facilityService: facilityService,
		viperConfig:     viperConfig,
	}
}

func (facilityHandler *Handler) FindAll(ginContext *gin.Context) {
	facilityResponses := facilityHandler.facilityService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Facility has been fetched", facilityResponses))
}

func (facilityHandler *Handler) FindAllPagination(ginContext *gin.Context) {
	paginationReq := model.PaginationRequest{
		Page:  1,
		Size:  10,
		Sort:  "id",
		Order: "asc",
	}

	// Bind query parameters to the struct
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	facilityResponses := facilityHandler.facilityService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Facility has been fetched", facilityResponses))
}

func (facilityHandler *Handler) CreateFacility(ginContext *gin.Context) {
	var createFacilityModel model.CreateFacilityRequest
	err := ginContext.ShouldBindBodyWithJSON(&createFacilityModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	facilityHandler.facilityService.Create(ginContext, &createFacilityModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Facility has been created", nil))
}

func (facilityHandler *Handler) UpdateFacility(ginContext *gin.Context) {
	var updateFacilityModel model.UpdateFacilityRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateFacilityModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	facilityHandler.facilityService.Update(ginContext, &updateFacilityModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Facility has been created", nil))
}

func (facilityHandler *Handler) DeleteFacility(ginContext *gin.Context) {
	var deleteBomModel model.DeleteFacilityRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteBomModel.Id = parsedBomId
	facilityHandler.facilityService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
