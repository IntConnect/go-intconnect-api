package facility

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/spf13/viper"
)

type Handler struct {
	facilityService Service
	formDecoder     *form.Decoder
	viperConfig     *viper.Viper
}

func NewHandler(facilityService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		facilityService: facilityService,
		viperConfig:     viperConfig,
		formDecoder:     form.NewDecoder(),
	}
}

func (facilityHandler *Handler) FindAllFacility(ginContext *gin.Context) {
	facilityResponses := facilityHandler.facilityService.FindAll()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Facility has been fetched", facilityResponses))
}

func (facilityHandler *Handler) FindAllFacilityPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := facilityHandler.facilityService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (facilityHandler *Handler) FindFacilityById(ginContext *gin.Context) {
	facilityId := ginContext.Param("id")
	parsedFacilityId, err := strconv.ParseUint(facilityId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	facilityResponse := facilityHandler.facilityService.FindById(parsedFacilityId)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("Facility has been fetched", facilityResponse))
}

func (facilityHandler *Handler) CreateFacility(ginContext *gin.Context) {
	var createFacilityModel model.CreateFacilityRequest
	err := ginContext.Request.ParseMultipartForm(32 << 20) // 32MB maxMemory
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	err = facilityHandler.formDecoder.Decode(&createFacilityModel, ginContext.Request.PostForm)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	thumbnailFile, err := ginContext.FormFile("thumbnail")
	createFacilityModel.Thumbnail = thumbnailFile
	facilityHandler.facilityService.Create(ginContext, &createFacilityModel)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse[interface{}]("Facility has been updated", nil))
}

func (facilityHandler *Handler) UpdateFacility(ginContext *gin.Context) {
	var updateFacilityModel model.UpdateFacilityRequest
	err := ginContext.Request.ParseMultipartForm(32 << 20) // 32MB maxMemory
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	err = facilityHandler.formDecoder.Decode(&updateFacilityModel, ginContext.Request.PostForm)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	thumbnailFile, _ := ginContext.FormFile("thumbnail")
	updateFacilityModel.Thumbnail = thumbnailFile
	facilityId := ginContext.Param("id")
	parsedFacilityId, err := strconv.ParseUint(facilityId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateFacilityModel.Id = parsedFacilityId
	facilityHandler.facilityService.Update(ginContext, &updateFacilityModel)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse[interface{}]("Facility has been updated", nil))
}

func (facilityHandler *Handler) DeleteFacility(ginContext *gin.Context) {
	var deleteFacilityModel model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteFacilityModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	facilityId := ginContext.Param("id")
	parsedFacilityId, err := strconv.ParseUint(facilityId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteFacilityModel.Id = parsedFacilityId
	paginatedResponse := facilityHandler.facilityService.Delete(ginContext, &deleteFacilityModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
