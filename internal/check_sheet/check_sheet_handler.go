package check_sheet

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
	checkSheetService Service
	viperConfig       *viper.Viper
}

func NewHandler(checkSheetService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		checkSheetService: checkSheetService,
		viperConfig:       viperConfig,
	}
}

func (checkSheetHandler *Handler) FindAllCheckSheet(ginContext *gin.Context) {
	checkSheetResponses := checkSheetHandler.checkSheetService.FindAll()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("Check sheet document template has been fetched", checkSheetResponses))
}

func (checkSheetHandler *Handler) FindAllCheckSheetPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := checkSheetHandler.checkSheetService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (checkSheetHandler *Handler) FindCheckSheetById(ginContext *gin.Context) {
	checkSheetId := ginContext.Param("id")
	parsedCheckSheetId, err := strconv.ParseUint(checkSheetId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := checkSheetHandler.checkSheetService.FindById(ginContext, parsedCheckSheetId)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (checkSheetHandler *Handler) CreateCheckSheet(ginContext *gin.Context) {
	var createCheckSheetModel model.CreateCheckSheetRequest
	err := ginContext.ShouldBindBodyWithJSON(&createCheckSheetModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedRes := checkSheetHandler.checkSheetService.Create(ginContext, &createCheckSheetModel)
	ginContext.JSON(http.StatusOK, paginatedRes)
}

func (checkSheetHandler *Handler) UpdateCheckSheet(ginContext *gin.Context) {
	var updateCheckSheetModel model.UpdateCheckSheetRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateCheckSheetModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrPayloadInvalid))
	checkSheetId := ginContext.Param("id")
	parsedCheckSheetId, err := strconv.ParseUint(checkSheetId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrParameterInvalid))
	updateCheckSheetModel.Id = parsedCheckSheetId
	paginatedRes := checkSheetHandler.checkSheetService.Update(ginContext, &updateCheckSheetModel)
	ginContext.JSON(http.StatusOK, paginatedRes)
}

func (checkSheetHandler *Handler) DeleteCheckSheet(ginContext *gin.Context) {
	var deleteResourceGeneralRequest model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteResourceGeneralRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrPayloadInvalid))
	checkSheetId := ginContext.Param("id")
	parsedCheckSheetId, err := strconv.ParseUint(checkSheetId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteResourceGeneralRequest.Id = parsedCheckSheetId
	paginatedRes := checkSheetHandler.checkSheetService.Delete(ginContext, &deleteResourceGeneralRequest)
	ginContext.JSON(http.StatusOK, paginatedRes)
}
