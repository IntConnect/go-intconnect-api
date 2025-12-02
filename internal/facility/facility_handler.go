package facility

import (
	"fmt"
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

func (facilityHandler *Handler) FindAll(ginContext *gin.Context) {
	facilityResponses := facilityHandler.facilityService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Facility has been fetched", facilityResponses))
}

func (facilityHandler *Handler) FindAllPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := facilityHandler.facilityService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (facilityHandler *Handler) CreateFacility(ginContext *gin.Context) {
	var createFacilityModel model.CreateFacilityRequest
	err := ginContext.Request.ParseMultipartForm(32 << 20) // 32MB maxMemory
	fmt.Println(err)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	err = facilityHandler.formDecoder.Decode(&createFacilityModel, ginContext.Request.PostForm)
	fmt.Println(2)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	thumbnailFile, err := ginContext.FormFile("thumbnail")
	fmt.Println(3, err)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	createFacilityModel.ThumbnailHeader = thumbnailFile

	paginatedResponse := facilityHandler.facilityService.Create(ginContext, &createFacilityModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (facilityHandler *Handler) UpdateFacility(ginContext *gin.Context) {
	var updateFacilityModel model.UpdateFacilityRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateFacilityModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	facilityHandler.facilityService.Update(ginContext, &updateFacilityModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Facility has been created", nil))
}

func (facilityHandler *Handler) DeleteFacility(ginContext *gin.Context) {
	var deleteFacilityModel model.DeleteResourceGeneralRequest
	facilityId := ginContext.Param("id")
	parsedFacilityId, err := strconv.ParseUint(facilityId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteFacilityModel.Id = parsedFacilityId
	facilityHandler.facilityService.Delete(ginContext, &deleteFacilityModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
